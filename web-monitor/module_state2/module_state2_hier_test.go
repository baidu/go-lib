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

// test normal cases for StateData.toHierStateData()
func TestToHierStateData_case0(t *testing.T) {
	sd := NewStateData()
	sd.SCounters.inc("baidu", 1)
	sd.States["state"] = "ok"
	sd.NumStates["num_state"] = 1

	hsd, err := toHierStateData(sd)
	if err != nil {
		t.Error("err in toHierStateData()")
	}

	if hsd.SCounters["baidu"] != int64(1) {
		t.Errorf("toHierStateData_case0(): Sounters[\"baidu\"] [%d] != 1", hsd.SCounters["baidu"])
	}

	if hsd.States["state"] != "ok" {
		t.Errorf("toHierStateData_case0(): States[\"state\"] value[%s] != ok",
			hsd.States["state"])
	}

	if hsd.NumStates["num_state"] != int64(1) {
		t.Errorf("toHierStateData_case0(): NumStates[\"num_state\"] value[%d] != 1",
			hsd.NumStates["num_state"])
	}
}

// test normal cases for StateData.toHierStateData()
func TestToHierStateData_case1(t *testing.T) {
	sd := NewStateData()
	sd.SCounters.inc("baidu.op.bfe", 1)

	_, err := toHierStateData(sd)
	if err != nil {
		t.Error("err in toHierStateData()")
	}
}
