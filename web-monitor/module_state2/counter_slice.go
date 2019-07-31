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

// get diff of two counters

/*
Usage:
    import "github.com/baidu/go-lib/web-monitor/module_state2"

    var counter module_state2.Counter
    var counterSlice *module_state2.CounterSlice
    var state *module_state2.State

    // usage 1: get diff once
    counterSlice.Set(counter)
    // make some update to counter here
    counterSlice.Set(counter)
    // get diff between update
    diff := counterSlice.Get()

    // usage 2: update diff periodically and get when needed
    var examCnt examCounter
    // update diff periodically
    counterSlice.Init(state, interval)
    // get diff
    diff := counterSlice.Get()
*/

package module_state2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/baidu/go-lib/web-monitor/web_params"
)

/* diff of two counters    */
type CounterSlice struct {
	lock sync.Mutex

	lastTime time.Time
	duration time.Duration

	countersLast Counters //  last absolute counter
	countersDiff Counters //  diff in last duration

	keyPrefix string     //  for key-value string
	programName   string //  program name, e.g., 'bfe', for displaying variable in key-value string
}

type CounterDiff struct {
	LastTime string // time till
	Duration int    // in second

	Diff Counters

	KeyPrefix string     // for key-value string
	ProgramName   string // for program name
}

// SetKeyPrefix sets key prefix
func (cs *CounterSlice) SetKeyPrefix(prefix string) {
	cs.keyPrefix = prefix
}

// SetProgramName sets program name
func (cs *CounterSlice) SetProgramName(programName string) {
	cs.programName = programName
}

// GetKeyPrefix gets key prefix
func (cs *CounterSlice) GetKeyPrefix() string {
	return cs.keyPrefix
}

// Set sets to counter slice
func (cs *CounterSlice) Set(counters Counters) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	if cs.countersLast == nil {
		// not initialized
		cs.lastTime = time.Now()
		cs.countersLast = counters.copy()
		cs.countersDiff = NewCounters()
	} else {
		now := time.Now()
		cs.duration = now.Sub(cs.lastTime)
		cs.lastTime = now

		cs.countersDiff = counters.diff(cs.countersLast)
		cs.countersLast = counters.copy()
	}
}

// Get gets diff from counter slice
func (cs *CounterSlice) Get() CounterDiff {
	var retVal CounterDiff

	cs.lock.Lock()
	defer cs.lock.Unlock()

	if cs.countersLast == nil {
		retVal.Diff = NewCounters()
	} else {
		retVal.LastTime = cs.lastTime.Format("2006-01-02 15:04:05")
		retVal.Duration = int(cs.duration.Seconds())
		retVal.Diff = cs.countersDiff.copy()
	}

	retVal.KeyPrefix = cs.keyPrefix
	retVal.ProgramName = cs.programName

	return retVal
}

// GetJson gets json format of counter diff
func (cs *CounterSlice) GetJson() ([]byte, error) {
	return json.Marshal(cs.Get())
}

func (cd CounterDiff) keyGen(str string, withProgramName bool) string {
	return KeyGen(str, cd.KeyPrefix, cd.ProgramName, withProgramName)
}

// KV outputs key-value string (lines of key:value) for CounterDiff
func (cd CounterDiff) KV() []byte {
	return cd.kv(false)
}

// KVWithProgramName outputs key-value string (lines of key:value) for CounterDiff, with program name
func (cd CounterDiff) KVWithProgramName() []byte {
	return cd.kv(true)
}

// kv outputs key-value string (lines of key:value) for CounterDiff
func (cd CounterDiff) kv(withProgramName bool) []byte {
	var buf bytes.Buffer

	for key, value := range cd.Diff {
		key = cd.keyGen(key, withProgramName)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	return buf.Bytes()
}

// FormatOutput formats output according format value in params
func (cd *CounterDiff) FormatOutput(params map[string][]string) ([]byte, error) {
	format, err := web_params.ParamsValueGet(params, "format")
	if err != nil {
		format = "json"
	}

	switch format {
	case "json":
		return json.Marshal(cd)
	case "hier_json":
		return GetCdHierJson(cd)
	case "kv":
		return cd.KV(), nil
	case "kv_with_program_name":
		return cd.KVWithProgramName(), nil
	default:
		return nil, fmt.Errorf("format not support: %s", format)
	}
}

// handleCounterSlice is go-routine for periodically get counter slice
func (cs *CounterSlice) handleCounterSlice(s *State, interval int) {
	for {
		counter := s.GetCounters()
		cs.Set(counter)

		leftSeconds := NextInterval(time.Now(), interval)
		time.Sleep(time.Duration(leftSeconds) * time.Second)
	}
}

// Init initializes the counter diff
// Params:
//    - s: module State
//    - interval: interval to compute between two counters
// Notice: use this method only when you need to get diff between two counters periodically
func (cs *CounterSlice) Init(s *State, interval int) {
	go cs.handleCounterSlice(s, interval)
}
