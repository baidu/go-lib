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
This program provides converting flat CounterDiff to hierarchical CounterDiff for json output

Usage:
    import "github.com/baidu/go-lib/web-monitor/module_state2"

    var cd module_state2.CounterDiff
    data, err := module_state2.GetCdHierJson(&cd)

*/

package module_state2

import (
	"encoding/json"
	"fmt"
)

// hierCounterDiff is hierarchical structure for counter diff
type hierCounterDiff struct {
	LastTime string // time till
	Duration int    // in second

	Diff hierCounters
	KeyPrefix string //  for key-value output
}

// toHierCounterDiff converts CounterDiff to hierCounterDiff
// Params:
//  - cd: flat counter diff
// Returns:
//  - *hierCounterDiff: hierarchical counter diff
//  - error: error msg
func toHierCounterDiff(cd *CounterDiff) (*hierCounterDiff, error) {
	var hcd hierCounterDiff
	var err error

	hcd.Diff, err = toHierCounters(cd.Diff)
	if err != nil {
		return nil, fmt.Errorf("toHierCounterDiff(): %s", err.Error())
	}

	hcd.LastTime = cd.LastTime
	hcd.Duration = cd.Duration
	hcd.KeyPrefix = cd.KeyPrefix

	return &hcd, nil
}

// GetCdHierJson gets hierarchical counter diff of json format
//Params:
//  - cd: flat counter diff
//Returns:
//  - []byte: json formated byte
//  - error: error msg
func GetCdHierJson(cd *CounterDiff) ([]byte, error) {
	hierCounterDiff, err := toHierCounterDiff(cd)
	if err != nil {
		return nil, fmt.Errorf("GetCdHierJson(): %s", err.Error())
	}

	return json.Marshal(hierCounterDiff)
}
