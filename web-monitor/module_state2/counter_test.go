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

func TestCounterIncDec(t *testing.T) {
	counters := NewCounters()
	counters.inc("test", 2)
	counters.dec("test", 1)

	copy := counters.copy()

	value, ok := copy["test"]
	if !ok || value != 1 {
		t.Error("Counters.inc() or Counters.dec() fail")
	}
}

func TestCounterInit(t *testing.T) {
	counters := NewCounters()

	keys := []string{"test1", "test2", "test3"}
	counters.init(keys)

	copy := counters.copy()

	for _, key := range keys {
		value, ok := copy[key]
		if !ok || value != 0 {
			t.Error("Counters.init() fail")
		}
	}
}

func TestCounterCopy(t *testing.T) {
	counters := NewCounters()
	counters["test"] = 123

	copy := counters.copy()

	value, ok := copy["test"]
	if !ok || value != 123 {
		t.Error("Counters.copy() fail")
	}
}

func TestCounterDiff_case1(t *testing.T) {
	counters := NewCounters()
	counters["test"] = 223

	last := NewCounters()
	last["test"] = 123

	diff := counters.diff(last)

	value, ok := diff["test"]
	if !ok || value != 100 {
		t.Error("Counters.diff() fail")
	}
}

func TestCounterDiff_case2(t *testing.T) {
	counters := NewCounters()
	counters["test"] = 123

	last := NewCounters()

	diff := counters.diff(last)

	value, ok := diff["test"]
	if !ok || value != 123 {
		t.Error("Counters.diff() fail")
	}
}

func TestCounterSum(t *testing.T) {
	counters1 := NewCounters()
	counters1["test1"] = 10
	counters1["test2"] = 20

	counters2 := NewCounters()
	counters2["test2"] = 20
	counters2["test3"] = 30

	counters1.Sum(counters2)
	if counters1["test1"] != 10 {
		t.Error("Counters.Sum() fail")
	}
	if counters1["test2"] != 40 {
		t.Error("Counters.Sum() fail")
	}
	if counters1["test3"] != 30 {
		t.Error("Counters.Sum() fail")
	}
}
