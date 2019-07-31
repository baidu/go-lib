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

package module_state2

import (
	"testing"
)

func TestNewNode_case0(t *testing.T) {
	key := "test"
	var value int64 = 100
	var tn *treeNode

	// isLastKey = true case
	tn = newNode(true, key, value)
	if tn.elem.key != key {
		t.Errorf("TestNewNode_case0(): key [%s] should be %s", tn.elem.key, key)
	}

	if tn.elem.value != value {
		t.Errorf("TestNewNode_case0(): value [%d] should be %d", tn.elem.value, value)
	}

	if tn.children != nil {
		t.Error("TestNewNode_case0(): children != nil")
	}

	// isLastKey = false case
	tn = newNode(false, key, value)
	if tn.elem.key != key {
		t.Errorf("TestNewNode_case0(): key [%s] should be %s", tn.elem.key, key)
	}

	if tn.elem.value != -1 {
		t.Errorf("TestNewNode_case0(): value [%d] should be %d", tn.elem.value, -1)
	}

	if tn.children != nil {
		t.Error("TestNewNode_case0(): children != nil")
	}
}

func TestCheckValidity_case0(t *testing.T) {
	var node treeNode
	var err error
	cntKey := "baidu.op.bfe"

	node.children = make([]*treeNode, 0)
	err = checkValidity(&node, true, cntKey)
	if err == nil {
		t.Error("TestCheckValidity_case0(): err should not = nil")
	}

	err = checkValidity(&node, false, cntKey)
	if err == nil {
		t.Error("TestCheckValidity_case0(): err should not = nil")
	}

	node.children = make([]*treeNode, 1)
	err = checkValidity(&node, true, cntKey)
	if err == nil {
		t.Error("TestCheckValidity_case0(): err should not = nil")
	}

	err = checkValidity(&node, false, cntKey)
	if err != nil {
		t.Error("TestCheckValidity_case0(): check should = nil")
	}
}

func TestGetNode_case0(t *testing.T) {
	var children []*treeNode
	var node *treeNode
	key := "bfe"

	node = getNode(children, key)
	if node != nil {
		t.Errorf("TestGetNode_case0(): node should = nil")
	}

	children = make([]*treeNode, 0)
	children = append(children, &treeNode{element{"root", -1}, nil})
	node = getNode(children, key)
	if node != nil {
		t.Errorf("TestGetNode_case0(): node should = nil")
	}

	children = append(children, &treeNode{element{"bfe", -1}, nil})
	node = getNode(children, key)
	if node == nil || node.elem.key != key {
		t.Errorf("TestGetNode_case0(): node should != nil or key[%s] != %s", node.elem.key, key)
	}

}

func TestInsert_case0(t *testing.T) {
	var father treeNode
	child := &treeNode{element{"root", -1}, nil}

	insert(&father, child)
	if len(father.children) != 1 || father.children[0].elem.key != "root" {
		t.Error("TestInsert_case0(): insert err")
	}
}

// normal cases
func TestNewMultiTree_case0(t *testing.T) {
	var err error
	c := NewCounters()
	c.inc("baidu.a.1", 1)
	_, err = newMultiTree(c)
	if err != nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}

	c.inc("baidu.a.2", 2)
	_, err = newMultiTree(c)
	if err != nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}

	c.inc("baidu.a.3", 3)
	_, err = newMultiTree(c)
	if err != nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}

	c.inc("baidu.b", 4)
	_, err = newMultiTree(c)
	if err != nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}

	c.inc("qq", 4)
	_, err = newMultiTree(c)
	if err != nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}
}

// abnormal case
func TestNewMultiTree_case1(t *testing.T) {
	var err error
	c := NewCounters()
	c.inc("baidu.a", 1)
	_, err = newMultiTree(c)
	if err != nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}

	c.inc("baidu.a.1", 2)
	_, err = newMultiTree(c)
	if err == nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}
}

// abnormal case
func TestNewMultiTree_case2(t *testing.T) {
	var err error
	c := NewCounters()
	c.inc("baidu.a.1", 2)
	_, err = newMultiTree(c)
	if err != nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}

	c.inc("baidu.a", 1)
	_, err = newMultiTree(c)
	if err == nil {
		t.Errorf("TestNewMultiTree_case0(): %s", err.Error())
	}
}
