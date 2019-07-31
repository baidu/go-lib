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

// collecting state info of a module
/*
Usage:
    import "github.com/baidu/go-lib/web-monitor/module_state2"

    var state module_state2.State

    state.Init()

    state.Inc("counter", 1)
    state.Set("state", "OK")
    state.SetNum("cap", 100)
    state.SetFloat("cap", 100.1)

    stateData := state.Get()
*/

package module_state2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

import (
	"github.com/baidu/go-lib/web-monitor/web_params"
)

// StateData holds state data
type StateData struct {
	SCounters     Counters          // for count up
	States        map[string]string // for store states
	NumStates     Counters          // for store num states
	FloatStates   FloatCounters     // for store float states
	KeyPrefix string                // for key-value output
	ProgramName   string            // for program name
}

// State is state with mutex protect
type State struct {
	lock sync.Mutex
	data StateData
}

func NewStateData() *StateData {
	sd := new(StateData)
	sd.SCounters = NewCounters()
	sd.States = make(map[string]string)
	sd.NumStates = NewCounters()
	sd.FloatStates = NewFloatCounters()

	return sd
}

// copy makes a copy for StateData
func (sd *StateData) copy() *StateData {
	copy := new(StateData)

	copy.SCounters = sd.SCounters.copy()

	copy.States = make(map[string]string)
	for key, value := range sd.States {
		copy.States[key] = value
	}

	copy.NumStates = NewCounters()
	for numKey, numValue := range sd.NumStates {
		copy.NumStates[numKey] = numValue
	}

	copy.FloatStates = NewFloatCounters()
	for floatKey, floatValue := range sd.FloatStates {
		copy.FloatStates[floatKey] = floatValue
	}

	copy.KeyPrefix = sd.KeyPrefix
	copy.ProgramName = sd.ProgramName

	return copy
}

func (sd *StateData) keyGen(key string, withProgramName bool) string {
	return KeyGen(key, sd.KeyPrefix, sd.ProgramName, withProgramName)
}

// KV output key-value string (lines of key:value) for StateData
func (sd *StateData) KV() []byte {
	return sd.kv(false)
}

// KVWithProgramName outputs key-value string (lines of key:value) for StateData, with program name
func (sd *StateData) KVWithProgramName() []byte {
	return sd.kv(true)
}

// escapeQuote escapes "\" 
func escapeQuote(value string) string {
	return strings.Replace(value, "\"", "\\\"", -1)
}

// kv outputs key-value string (lines of key:value) for StateData
func (sd *StateData) kv(withProgramName bool) []byte {
	var buf bytes.Buffer

	// print SCounters
	for key, value := range sd.SCounters {
		key = sd.keyGen(key, withProgramName)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	// print States
	for key, value := range sd.States {
		key = sd.keyGen(key, withProgramName)
		value = escapeQuote(value)
		str := fmt.Sprintf("%s:\"%s\"\n", key, value)
		buf.WriteString(str)
	}

	// print NumStates
	for key, value := range sd.NumStates {
		key = sd.keyGen(key, withProgramName)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	// print floatStates
	for key, value := range sd.FloatStates {
		key = sd.keyGen(key, withProgramName)
		str := fmt.Sprintf("%s:%f\n", key, value)
		buf.WriteString(str)
	}

	return buf.Bytes()
}

// FormatOutput formats output according format value in params
func (sd *StateData) FormatOutput(params map[string][]string) ([]byte, error) {
	format, err := web_params.ParamsValueGet(params, "format")
	if err != nil {
		format = "json"
	}

	switch format {
	case "json":
		return json.Marshal(sd)
	case "hier_json":
		return GetSdHierJson(sd)
	case "kv":
		return sd.KV(), nil
	case "kv_with_program_name":
		return sd.KVWithProgramName(), nil
	default:
		return nil, fmt.Errorf("format not support: %s", format)
	}
}

// Init initializes the state
func (s *State) Init() {
	s.data.SCounters = NewCounters()
	s.data.States = make(map[string]string)
	s.data.NumStates = NewCounters()
	s.data.FloatStates = NewFloatCounters()
}

// SetKeyPrefix sets key prefix
func (s *State) SetKeyPrefix(prefix string) {
	s.data.KeyPrefix = prefix
}

// SetProgramName sets program name
func (s *State) SetProgramName(programName string) {
	s.data.ProgramName = programName
}

// Inc increases value for given key
func (s *State) Inc(key string, value int) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.SCounters.inc(key, value)
	s.lock.Unlock()
}

// Dec decreases value for given key
func (s *State) Dec(key string, value int) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.SCounters.dec(key, value)
	s.lock.Unlock()
}

// CountersInit Initializes counters for given keys to zero
func (s *State) CountersInit(keys []string) {
	s.lock.Lock()
	s.data.SCounters.init(keys)
	s.lock.Unlock()
}

// Set sets value to given key
func (s *State) Set(key string, value string) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.States[key] = value
	s.lock.Unlock()
}

// Delete deletes state for given key
func (s *State) Delete(key string) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	delete(s.data.States, key)
	s.lock.Unlock()
}

// SetNum sets num state to given key
func (s *State) SetNum(key string, value int64) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.NumStates[key] = value
	s.lock.Unlock()
}

// SetFloat sets float state to given key
func (s *State) SetFloat(key string, value float64) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.FloatStates[key] = value
	s.lock.Unlock()
}

// GetCounter gets counter value of given key
func (s *State) GetCounter(key string) int64 {
	s.lock.Lock()
	value, ok := s.data.SCounters[key]
	s.lock.Unlock()

	if !ok {
		value = 0
	}

	return value
}

// GetCounters gets all counters
func (s *State) GetCounters() Counters {
	s.lock.Lock()
	counters := s.data.SCounters.copy()
	s.lock.Unlock()

	return counters
}

// GetState get state value of given key
func (s *State) GetState(key string) string {
	s.lock.Lock()
	value, ok := s.data.States[key]
	s.lock.Unlock()

	if !ok {
		value = ""
	}

	return value
}

// GetNumState gets num state value of given key
func (s *State) GetNumState(key string) int64 {
	s.lock.Lock()
	value, ok := s.data.NumStates[key]
	s.lock.Unlock()

	if !ok {
		value = 0
	}

	return value
}

// GetFloatState gets float state value of given key
func (s *State) GetFloatState(key string) float64 {
	s.lock.Lock()
	value, ok := s.data.FloatStates[key]
	s.lock.Unlock()

	if !ok {
		value = float64(0.0)
	}

	return value
}

// GetAll gets all states
func (s *State) GetAll() *StateData {
	s.lock.Lock()
	copy := s.data.copy()
	s.lock.Unlock()
	return copy
}

// GetKeyPrefix gets key prefix
func (s *State) GetKeyPrefix() string {
	return s.data.KeyPrefix
}
