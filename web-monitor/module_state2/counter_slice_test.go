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
	"fmt"
	"testing"
	"time"
)

func TestCounterSliceGet(t *testing.T) {
	var cs CounterSlice

	diff := cs.Get()
	if len(diff.Diff) != 0 {
		t.Error("data in diff should be zero")
	}

	counters := NewCounters()
	counters["test"] = 123

	cs.Set(counters)
	diff = cs.Get()
	if len(diff.Diff) != 0 {
		t.Error("data in diff should be zero")
	}

	time.Sleep(time.Second)

	counters["test"] = 223
	cs.Set(counters)
	diff = cs.Get()

	value, ok := diff.Diff["test"]
	if !ok || value != 100 {
		t.Error("the diff for test should be 100")
	}
	if diff.Duration != 1 {
		t.Error("duration should be 1")
	}
	fmt.Printf("diff=%v\n", diff)
}

func TestCounterDiff_KV(t *testing.T) {
	var diff CounterDiff

	// prepare data
	diff.LastTime = "1234"
	diff.Duration = 5678

	diff.Diff = NewCounters()
	diff.Diff.inc("counter", 1)

	// output key-value string
	strOK := "counter:1\n"
	if string(diff.KV()) != strOK {
		t.Error("err in CounterDiff.KV()")
	}
}

func TestCounterDiff_KV_with_progName(t *testing.T) {
	var diff CounterDiff

	// prepare data
	diff.LastTime = "1234"
	diff.Duration = 5678

	diff.Diff = NewCounters()
	diff.Diff.inc("counter", 1)
	diff.ProgramName = "program"

	// output key-value string
	strOK := "program.counter:1\n"
	if string(diff.KVWithProgramName()) != strOK {
		t.Error("err in CounterDiff.KVWithProgramName()")
	}
}

func TestFormatOutput4CounterSlice(t *testing.T) {
	var err error
	var cs CounterSlice
	var diff CounterDiff

	sd := NewStateData()
	sd.SCounters.init([]string{
		"baidu",
	})

	diff = cs.Get()
	_, err = (&diff).FormatOutput(map[string][]string{"param": []string{
		"no_format"},
	})

	if err != nil {
		t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
	}

	diff = cs.Get()
	_, err = (&diff).FormatOutput(map[string][]string{"format": []string{
		"json"},
	})

	if err != nil {
		t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
	}

	diff = cs.Get()
	_, err = (&diff).FormatOutput(map[string][]string{"format": []string{
		"hier_json"},
	})

	if err != nil {
		t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
	}

	diff = cs.Get()
	_, err = (&diff).FormatOutput(map[string][]string{"format": []string{
		"kv"},
	})

	if err != nil {
		t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
	}

	diff = cs.Get()
	_, err = (&diff).FormatOutput(map[string][]string{"format": []string{
		"no_exist"},
	})

	if err == nil {
		t.Error("TestFormatOutCD_Case0(): err should not equal nil")
	}
}
