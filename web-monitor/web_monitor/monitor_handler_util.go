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

// utility functions for simplifying use of the library

package web_monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
)

import (
	"github.com/baidu/go-lib/web-monitor/delay_counter"
	"github.com/baidu/go-lib/web-monitor/module_state2"
	"github.com/baidu/go-lib/web-monitor/kv_encode"
	"github.com/baidu/go-lib/web-monitor/web_params"
)

// function prototype for getting CounterDiff/DelayOutput/StateData
type GetCounterDiffFunc func() *module_state2.CounterDiff
type GetDelayOutputFunc func() *delay_counter.DelayOutput
type GetStateDataFunc func() *module_state2.StateData

// CreateStateDataHandler creates monitor handler for StateData
//
// Params:
//     - getter: func for getting StateData
//
// Returns:
//     - a monitor handler
func CreateStateDataHandler(getter GetStateDataFunc) interface{} {
	return func(params map[string][]string) ([]byte, error) {
		var buff []byte
		var err error

		// get StateData
		state := getter()
		if state == nil {
			return nil, fmt.Errorf("GetStateDataFunc: invalid data")
		}

		// return encoded data
		format := GetFormatParam(params)
		switch format {
		case "json":
			buff, err = json.Marshal(state)
		case "kv":
			buff = state.KV()
		case "kv_with_program_name":
			buff = state.KVWithProgramName()
		default:
			err = fmt.Errorf("invalid format:%s", format)
		}
		return buff, err
	}
}

// CreateDelayOutputHandler creates monitor handler for DelayRecent
//
// Params:
//     - getter: func for getting DelayOutput
//
// Returns:
//     - a monitor handler
func CreateDelayOutputHandler(getter GetDelayOutputFunc) interface{} {
	return func(params map[string][]string) ([]byte, error) {
		var buff []byte
		var err error

		// get DelayOutput
		delay := getter()
		if delay == nil {
			return nil, fmt.Errorf("GetDelayOutputFunc: invalid data")
		}

		// return encoded data
		format := GetFormatParam(params)
		switch format {
		case "json":
			buff, err = delay.GetJson()
		case "kv":
			buff = delay.GetKV()
		case "kv_with_program_name":
			buff = delay.GetKVWithProgramName()
		default:
			err = fmt.Errorf("invalid format:%s", format)
		}
		return buff, err
	}
}

// CreateCounterDiffHandler creates monitor handler for CounterDiff
//
// Params:
//     - getter: func for getting CounterDiff
//
// Returns:
//     - a monitor handler
func CreateCounterDiffHandler(getter GetCounterDiffFunc) interface{} {
	return func(params map[string][]string) ([]byte, error) {
		var buff []byte
		var err error

		// get CounterDiff
		diff := getter()
		if diff == nil {
			return nil, fmt.Errorf("GetCounterDiffFunc: invalid data")
		}

		// return encoded data
		format := GetFormatParam(params)
		switch format {
		case "json":
			buff, err = json.Marshal(diff)
		case "kv":
			buff = diff.KV()
		case "kv_with_program_name":
			buff = diff.KVWithProgramName()
		default:
			err = fmt.Errorf("invalid format:%s", format)
		}
		return buff, err
	}
}

// CreateMemStatsHandler creates monitor handler for getting memory statistics
//
// Params:
//     - keyPrefix: prefix of key, eg. <ServerName>_mem_stats
//
// Return:
//     - a monitor handler
func CreateMemStatsHandler(keyPrefix string) interface{} {
	return func(params map[string][]string) ([]byte, error) {
		var buff []byte
		var err error

		// get memory statistics
		var stat runtime.MemStats
		runtime.ReadMemStats(&stat)

		// return encoded data
		format := GetFormatParam(params)
		switch format {
		case "json":
			buff, err = json.Marshal(stat)
		case "kv":
			buff, err = MemStatsKVEncode(stat, keyPrefix)
		default:
			err = fmt.Errorf("invalid format:%s", format)
		}
		return buff, err
	}
}

// MemStatsKVEncode gets encode data of MemStats in key-value format
func MemStatsKVEncode(stat runtime.MemStats, keyPrefix string) ([]byte, error) {
	// for fields of baisc type
	buff, err := kv_encode.EncodeData(stat, keyPrefix, true)
	if err != nil {
		return nil, err
	}

	// for special fields
	prefix := keyPrefix
	if prefix != "" {
		prefix = prefix + "_"
	}
	var data bytes.Buffer

	// Note: stat.PauseNs is circular buffer of recent GC pause durations,
	// most recent at [(NumGC+255) % 256]
	data.WriteString(fmt.Sprintf("%s%s:%d\n", prefix, "LastPauseNs",
		stat.PauseNs[(stat.NumGC+255)%256]))

	buff = append(buff, data.Bytes()...)
	return buff, nil
}

// get format parameter
func GetFormatParam(params map[string][]string) string {
	format, err := web_params.ParamsValueGet(params, "format")
	if err != nil {
		format = "json" // default format is json
	}
	return format
}
