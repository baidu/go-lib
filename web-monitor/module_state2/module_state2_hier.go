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
This program provides converting flat StateDate to hierarchical StateData for json output

Usage:
    import "github.com/baidu/go-lib/web-monitor/module_state2"

    var sd module_state2.StateData
    data, err := GetSdHierJson(&sd)
*/

package module_state2

import (
	"encoding/json"
	"fmt"
)

// hierStateData holds hierarchical structure for StateData
type hierStateData struct {
	SCounters     hierCounters      // for count up
	States        map[string]string // for store states
	NumStates     hierCounters      // for store num states
	KeyPrefix string
}

// toHierStateData converts StateData to hierStateData
// Params:
//  - sd: flat state data
// Returns:
//  - *hierStateData: hierarchical state data
//  - error: error msg
func toHierStateData(sd *StateData) (*hierStateData, error) {
	var hsd hierStateData
	var err error

	hsd.SCounters, err = toHierCounters(sd.SCounters)
	if err != nil {
		return nil, fmt.Errorf("toHierStateData(): Scounters %s", err.Error())
	}

	hsd.States = sd.States
	hsd.NumStates, err = toHierCounters(sd.NumStates)
	if err != nil {
		return nil, fmt.Errorf("toHierStateData(): NumStates %s", err.Error())
	}

	hsd.KeyPrefix = sd.KeyPrefix

	return &hsd, nil
}

// GetSdHierJson gets hierarchical StataData of json format
// Params:
//  - sd: flat state data
// Returns:
//  - []byte: json formated byte
//  - error: error msg
func GetSdHierJson(sd *StateData) ([]byte, error) {
	hierState, err := toHierStateData(sd)
	if err != nil {
		return nil, fmt.Errorf("GetSdHierJson(): %s", err.Error())
	}

	return json.Marshal(hierState)
}
