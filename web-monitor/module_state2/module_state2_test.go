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
)

func TestModuleState(t *testing.T) {
	var state State
	var ok bool
	var value int64
	var vStr string
	var num int64

	state.Init()
	state.Inc("counter", 1)
	state.Inc("counter", 2)
	state.Dec("counter", 1)
	state.Set("state", "OK")
	state.SetNum("cap", 100)

	// test GetAll()
	data := state.GetAll()
	fmt.Println(*data)

	value, ok = data.SCounters["counter"]
	if !ok || value != 2 {
		t.Error("err in GetAll(), value should be 2")
	}

	vStr, ok = data.States["state"]
	if !ok || vStr != "OK" {
		t.Error("err in GetAll(), value should be OK")
	}

	// test GetCounter()
	value = state.GetCounter("counter")
	if value != 2 {
		t.Error("err in GetCounter(), value should be 2")
	}

	// test GetCounters()
	state.Inc("counter2", 3)
	counters := state.GetCounters()
	value, ok = counters["counter"]
	if !ok || value != 2 {
		t.Error("err in GetCounters(), value should be 2")
	}
	value, ok = counters["counter2"]
	if !ok || value != 3 {
		t.Error("err in GetCounters(), value should be 3")
	}

	// test GetState()
	vStr = state.GetState("state")
	if vStr != "OK" {
		t.Error("err in GetState(), value should be OK")
	}

	// test GetNumState()
	num = state.GetNumState("cap")
	if num != 100 {
		t.Error("err in GetNumSate(), num should be 100")
	}
}

func TestModuleStateCountersInit(t *testing.T) {
	var state State

	state.Init()

	// init counters
	keys := []string{"test1", "test2", "test3"}
	state.CountersInit(keys)

	// check counters
	counters := state.GetCounters()
	for _, key := range keys {
		value, ok := counters[key]

		if !ok || value != 0 {
			t.Error("err in CountersInit(), value should be 0")
		}
	}
}

func TestModuleStateNil(t *testing.T) {
	// test support of nil for Inc(), Dec(), Set(), SetNum()
	var pState *State

	if pState != nil {
		t.Error("pState should be nil")
	}

	pState.Inc("test", 1)
	pState.Dec("test", 1)
	pState.Set("state", "ok")
	pState.SetNum("num", 1)
}

// test for StateData.KV()
func TestStateData_KV(t *testing.T) {
	sd := NewStateData()

	sd.SCounters.inc("counter", 1)
	sd.States["state"] = "ok"
	sd.NumStates["num_state"] = 1

	strOK := "counter:1\n" + "state:\"ok\"\n" + "num_state:1\n"

	if string(sd.KV()) != strOK {
		t.Error("err in StateData.KV()")
	}

	sd = NewStateData()
	sd.SCounters.inc("TLS_ALPN_SPDY/3.1", 1)
	sd.States["TLS_ALPN_SPDY/2.1"] = "ok"
	sd.NumStates["TLS_ALPN_SPDY/1.1"] = 1

	strOK = "TLS_ALPN_SPDY_3.1:1\n" + "TLS_ALPN_SPDY_2.1:\"ok\"\n" + "TLS_ALPN_SPDY_1.1:1\n"

	if string(sd.KV()) != strOK {
		t.Error("err in StateData.KV()")
	}

}

func TestFormatOutput4StateData(t *testing.T) {
	var err error
	var b []byte
	var result string

	s := NewStateData()
	s.SCounters.init([]string{
		"baidu",
	})

	_, err = s.FormatOutput(map[string][]string{"param": []string{
		"no_format"},
	})

	if err != nil {
		t.Errorf("TestFormatOutSD_Case0(): %s", err.Error())
	}

	b, err = s.FormatOutput(map[string][]string{"format": []string{
		"json"},
	})

	if err != nil {
		t.Errorf("TestFormatOutSD_Case0(): %s", err.Error())
	}

	result = "{\"SCounters\":{\"baidu\":0},\"States\":{},\"NumStates\":{}," +
		"\"FloatStates\":{},\"KeyPrefix\":\"\",\"ProgramName\":\"\"}"

	if string(b) != result {
		t.Errorf("TestFormatOutSD_Case0(): %s not equal %s", string(b), result)
	}

	_, err = s.FormatOutput(map[string][]string{"format": []string{
		"hier_json"},
	})

	if err != nil {
		t.Errorf("TestFormatOutSD_Case0(): %s", err.Error())
	}

	if string(b) != result {
		t.Errorf("TestFromatOutSD_Case0(): %s not equal %s", string(b), result)
	}

	b, err = s.FormatOutput(map[string][]string{"format": []string{
		"kv"},
	})

	if err != nil {
		t.Errorf("TestFormatOutSD_Case0(): %s", err.Error())
	}

	if string(b) != "baidu:0\n" {
		t.Errorf("TestFormatOutSD_Case0(): %s not equal baidu:0", string(b))
	}

	_, err = s.FormatOutput(map[string][]string{"format": []string{
		"no_exist"},
	})

	if err == nil {
		t.Error("TestServiceCounterGet_Case0(): err should not equal nil")
	}
}
