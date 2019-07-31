// Copyright (c) 2018 Baidu, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
   This program converts Counters(flat display counters) to
   hierCounters(hierarchical display counters) by building
   a multi tree.
   Convert happens only when key in Counters is splited by dot(".")
   -------------------------------------------------------------------
   For example:

                                            "CheckOnly"           : 1
                                            "NoCheck"             : 2
   "CheckOnly"                  : 4,        "WaitResponse : {
   "NoCheck"                    : 6,            "Forbidden"       : 4
   "WaitResponse.Forbidden"     : 8,  --->          "Pass" : {
   "WaitResponse.Pass.OK"       : 10                    "OK"      : 6
   "WaitResponse.Pass.Timeout"  : 10                    "Timeout" : 8
   "WaitResponse.Pass.Other"    : 10                    "Other"   :10
                                                }
                                            }
  --------------------------------------------------------------------
*/

package module_state2

import (
	"fmt"
	"strings"
)

// a multiTree is composed by many nodes: root node and data nodes.
// root node is just used to build the tree, data nodes contain the meta data
// e.g.
// multiTree after build "WaitResponse.Forbidden :10" is as follows:
//                            "root:-1"
//                                |
//                          "WaitResponse:-1"
//                                |
//                           "Forbidden:10"

// element is element of tree node, which contains the meta data
type element struct {
	// key is name of the node. for root node, key is "root"
	key string
	
	// value is the value of the node
	//  - for non-leaf node, value is always -1
	//  - for leaf node, value is the counter
	value int64
}

// treeNode is struct of tree node
type treeNode struct {
	// elem is element of the node
	elem     element

	// children are children of the node
	children []*treeNode
}

// newNode creates a new node with key and value
// Params:
//  - isLastKey: nodeKey is/isn't the last key of the dot splited cntKey
//  - nodeKey: key of the node
//  - nodeVal: value of the counter key
// Returns:
//  - pointer to the new node
func newNode(isLastKey bool, nodeKey string, nodeVal int64) *treeNode {
	var node *treeNode

	if isLastKey {
		// for leaf node, value is set to the counter
		node = &treeNode{element{nodeKey, nodeVal}, nil}
	} else {
		// for non-leaf node, value is set to -1
		node = &treeNode{element{nodeKey, -1}, nil}
	}

	return node
}

// checkValidity check if it ok to build the tree or not
// Params:
//  - node: the checked node
//  - isLastKey: nodeKey is/isn't the last key of the dot splited cntKey
//  - cntKey: counter key
// Returns: error
//  - nil if check ok
//  - err info if check fail
func checkValidity(node *treeNode, isLastKey bool, cntKey string) error {
	isLeafNode := len(node.children) == 0

	if isLeafNode {
		return fmt.Errorf("key[%s] at least has one prefix key", cntKey)
	} else {
		if isLastKey {
			return fmt.Errorf("key[%s] is the prefix of other keys", cntKey)
		} else {
			return nil
		}
	}
}

// getNode gets node among the children whose element.key == key
// Params:
//  - children: node set
//  - key: key string
// Returns: *treeNode
// - node contains the key when get
// - nil while no get
func getNode(children []*treeNode, key string) *treeNode {
	for i := 0; i < len(children); i++ {
		if children[i].elem.key == key {
			return children[i]
		}
	}

	return nil
}

// insert inserts a node as child to another node
// Params:
//  - father: father node
//  - child: child node
func insert(father *treeNode, child *treeNode) {
	if father.children == nil {
		father.children = make([]*treeNode, 0)
	}

	father.children = append(father.children, child)
}

// buildTree builds multiTree from specified Counter key
// Params:
//  - t: root node of the built tree
//  - cntKey: counter key, which is a dot splited string, e.g., a.b.c
//  - cntVal: value related to the counter key
func buildTree(t *treeNode, cntKey string, cntVal int64) error {
	// use dot to separate cntrKey, each key in the slice stands for a tree node
	// e.g., a.b.c => root->node(a)->node(b)->node(c)
	keySlice := strings.Split(cntKey, ".")

	for i := 0; i < len(keySlice); i++ {
		isLastKey := (i + 1) == len(keySlice)
		nodeKey := keySlice[i]

		// get children node of t which key is nodeKey
		node := getNode(t.children, nodeKey)
		if node == nil {
			// node does not exist, create a new node
			node = newNode(isLastKey, nodeKey, cntVal)
			// insert the new node as a child of t
			insert(t, node)
		} else {
			// node exist, check validity
			err := checkValidity(node, isLastKey, cntKey)
			if err != nil {
				return err
			} else {
				//bypass, do nothing
			}
		}

		// to build children of node
		t = node
	}
	return nil
}

// newMultiTree creates multiTree from Counters
// Params:
//  -c: the counters used to build multiTree
// Returns: (*treeNode, error)
//  - *treeNode: root node of the built Tree
//  - error: err info
func newMultiTree(c Counters) (*treeNode, error) {
	// new root node
	root := &treeNode{element{"root", -1}, nil}

	for cntKey, cntVal := range c {
		// build tree with every key in counter
		err := buildTree(root, cntKey, cntVal)
		if err != nil {
			return nil, err
		}
	}

	return root, nil
}
