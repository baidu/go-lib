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

func TestToHierCounterDiff_case0(t *testing.T) {
	var cd CounterDiff
	cd.LastTime = "lastTime"
	cd.Duration = 20
	cd.Diff = NewCounters()
	cd.Diff.inc("baidu.op", 1)
	hcd, err := toHierCounterDiff(&cd)
	if err != nil {
		t.Errorf("TestToHierCounterDiff(): %s", err.Error())
	}

	if hcd.LastTime != cd.LastTime {
		t.Errorf("hcd.LastTime[%s] != cd.LastTime[%s]", hcd.LastTime, cd.LastTime)
	}

	if hcd.Duration != cd.Duration {
		t.Errorf("hcd.Duration[%d] != cd.Duration[%d]", hcd.Duration, cd.Duration)
	}

}
