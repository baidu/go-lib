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

// init case
func TestInit_case0(t *testing.T) {
	c := NewCounters()
	c.inc("baidu.a.1", 1)
	c.inc("baidu.a.2", 2)
	c.inc("baidu.b", 3)

	mt, _ := newMultiTree(c)

	hc := newHierCounters()
	hc.init(mt)

	// for baidu.b
	if v1, ok1 := hc["baidu"].(hierCounters); ok1 {
		if v2, ok2 := v1["b"].(int64); ok2 {
			if v2 != int64(3) {
				t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"b\"] (%d) != 3",
					((hc["baidu"].(hierCounters))["b"]).(int64))
			}
		} else {
			t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"b\"] is not an int64")
		}
	} else {
		t.Errorf("TestToHierCounters(): hc[\"baidu\"] is not an Counters")
	}

	// for baidu.a.1
	if v3, ok3 := hc["baidu"].(hierCounters); ok3 {
		if v4, ok4 := v3["a"].(hierCounters); ok4 {
			if v5, ok5 := v4["1"].(int64); ok5 {
				if v5 != int64(1) {
					t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"a\"][\"1\"] (%d) != 1",
						(((hc["baidu"].(hierCounters))["a"]).(hierCounters))["1"].(int64))
				}
			} else {
				t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"a\"][\"1\"] is not an int64")
			}
		} else {
			t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"a\"] is not an Counters")
		}
	} else {
		t.Errorf("TestToHierCounters(): hc[\"baidu\"] is not an Counters")
	}

	// for baidu.a.2
	if v3, ok3 := hc["baidu"].(hierCounters); ok3 {
		if v4, ok4 := v3["a"].(hierCounters); ok4 {
			if v5, ok5 := v4["2"].(int64); ok5 {
				if v5 != int64(2) {
					t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"a\"][\"2\"] (%d) != 2",
						(((hc["baidu"].(hierCounters))["a"]).(hierCounters))["2"].(int64))
				}
			} else {
				t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"a\"][\"2\"] is not an int64")
			}
		} else {
			t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"a\"] is not an Counters")
		}
	} else {
		t.Errorf("TestToHierCounters(): hc[\"baidu\"] is not an Counters")
	}
}

// normal cases
func TestToHierCounters_case0(t *testing.T) {
	c := NewCounters()
	c.inc("baidu", 1)
	_, err := toHierCounters(c)
	if err != nil {
		t.Errorf("TestToHierCounters(): %s", err.Error())
	}
}

// normal cases
func TestToHierCounters_case1(t *testing.T) {
	c := NewCounters()
	c.inc("baidu", 1)
	c.inc("baidu.a", 1)
	_, err := toHierCounters(c)
	if err == nil {
		t.Error("TestToHierCounters(): err must not be nil")
	}
}

// normal cases
func TestToHierCounters_case2(t *testing.T) {
	c := NewCounters()
	hc, err := toHierCounters(c)
	if err != nil {
		t.Errorf("TestToHierCounters(): %s", err.Error())
	}

	if len(hc) != 0 {
		t.Errorf("TestToHierCounters(): len(hc)[%d] != 0", len(hc))
	}
}
