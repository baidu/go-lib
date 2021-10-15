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

package delay_counter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/baidu/go-lib/web-monitor/module_state2"
	"github.com/baidu/go-lib/web-monitor/web_params"
)

type DelayRecent struct {
	lock sync.Mutex

	interval int // interval of making switch

	currTime time.Time
	current  DelaySummary // data for current minute

	pastTime time.Time
	past     DelaySummary // data for last minute

	// for key-value output
	KeyPrefix   string // prefix for key
	ProgramName string // program name
}

// DelayOutput is designed for json output
type DelayOutput struct {
	Interval    int
	KeyPrefix   string
	ProgramName string

	CurrTime string
	Current  DelaySummary

	PastTime string
	Past     DelaySummary
}

// Init initializes delay table
//
// Params:
//      - interval: interval for move current to past
//      - bucketSize: size of each delay bucket, e.g., 1(ms) or 2(ms)
//      - number of bucket
//
func (t *DelayRecent) Init(interval int, bucketSize int, bucketNum int) {
	t.currTime = time.Now()
	// adjust time
	t.currTime = t.currTime.Truncate(time.Duration(interval) * time.Second)

	t.interval = interval

	// initialize DelayCounters
	t.current.Init(bucketSize, bucketNum)
	t.past.Init(bucketSize, bucketNum)
}

// SetKeyPrefix sets prefix used in key generation
func (t *DelayRecent) SetKeyPrefix(prefix string) {
	t.KeyPrefix = prefix
}

// SetProgramName sets program name used in key generation
func (t *DelayRecent) SetProgramName(programName string) {
	t.ProgramName = programName
}

// AddBySub adds one new data to the table, by providing start time and end time
func (t *DelayRecent) AddBySub(start time.Time, end time.Time) {
	/* get duration from start to now, in Microsecond   */
	duration := end.Sub(start).Nanoseconds() / 1000

	t.Add(duration)
}

// Clear clears counters
func (t *DelayRecent) Clear() {
	t.current.Clear()
	t.past.Clear()
}

// Add adds one new data to the table.
//
// Params:
//     - duration: delay duration, in Microsecond (10^-6)
func (t *DelayRecent) Add(duration int64) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.trySwitch()
	t.current.Add(duration)
}

// AddDuration adds one new data to the table.
//
// Params:
//      - duration: time duration of delay (in Nanosecond)
func (t *DelayRecent) AddDuration(duration time.Duration) {
	delay := int64(duration / time.Microsecond)
	t.Add(delay)
}

// trySwitch checks and switches DelayRecent
func (t *DelayRecent) trySwitch() {
	now := time.Now()
	if (t.currTime.Unix() / int64(t.interval)) != (now.Unix() / int64(t.interval)) {
		/* they are not in the same minute, do a switch */
		t.pastTime = t.currTime
		t.currTime = now

		t.past.Copy(t.current)
		t.current.Clear() // clear t.current
	}
}

func (t *DelayRecent) get() DelayOutput {
	var retVal DelayOutput

	t.lock.Lock()
	defer t.lock.Unlock()

	t.trySwitch()
	retVal.Interval = t.interval
	retVal.CurrTime = fmt.Sprintf(t.currTime.Format("2006-01-02 15:04:05"))
	retVal.Current.Copy(t.current)
	retVal.PastTime = fmt.Sprintf(t.pastTime.Format("2006-01-02 15:04:05"))
	retVal.Past.Copy(t.past)

	// set key prefix and program name
	retVal.KeyPrefix = t.KeyPrefix
	retVal.ProgramName = t.ProgramName

	return retVal
}

// Get gets counter from table
func (t *DelayRecent) Get() DelayOutput {
	retVal := t.get()

	// calc average
	retVal.Current.CalcAvg()
	retVal.Past.CalcAvg()

	return retVal
}

// GetJson gets data in the table, return with json string
func (t *DelayRecent) GetJson() ([]byte, error) {
	d := t.Get()
	return d.GetJson()
}

// get data in the table, return with key-value string (i.e., lines of key:value)
func (t *DelayRecent) GetKV() []byte {
	d := t.Get()
	return d.GetKV()
}

// get data in the table, return with key-value string, with program name
func (t *DelayRecent) GetKVWithProgramName() []byte {
	d := t.Get()
	return d.GetKVWithProgramName()
}

// get data in the table, return with prometheus format
func (t *DelayRecent) GetPrometheusFormat() []byte {
	d := t.Get()
	return d.GetPrometheusFormat()
}

// FormatOutput formats output according to format value in params
func (t *DelayRecent) FormatOutput(params map[string][]string) ([]byte, error) {
	format, err := web_params.ParamsValueGet(params, "format")
	if err != nil {
		format = "json"
	}

	switch format {
	case "json", "hier_json":
		return t.GetJson()
	case "kv", "noah":
		return t.GetKV(), nil
	case "kv_with_program_name":
		return t.GetKVWithProgramName(), nil
	case "prometheus":
		return t.GetPrometheusFormat(), nil
	default:
		return nil, fmt.Errorf("format not support: %s", format)
	}
}

// Sum calculates sum of DelayOutput
func (d *DelayOutput) Sum(d2 DelayOutput) error {
	if d.Interval != d2.Interval {
		return fmt.Errorf("Interval not match")
	}

	if err := d.Current.calcSum(d2.Current); err != nil {
		return err
	}
	if err := d.Past.calcSum(d2.Past); err != nil {
		return err
	}

	if d.CurrTime < d2.CurrTime {
		d.CurrTime = d2.CurrTime
	}
	if d.PastTime < d2.PastTime {
		d.PastTime = d2.PastTime
	}
	return nil
}

// get json string for DelayOutput
func (d *DelayOutput) GetJson() ([]byte, error) {
	return json.Marshal(d)
}

// generate key prefix
func (d *DelayOutput) keyPrefixGen(key string, withProgramName bool) string {
	return module_state2.KeyGen(key, d.KeyPrefix, d.ProgramName, withProgramName)
}

// get key-value string for DelayOutput, without program name
func (d *DelayOutput) GetKV() []byte {
	return d.getKV(false)
}

// GetKVStringWithProgramName gets key-value string for DelayOutput, with program name
func (d *DelayOutput) GetKVWithProgramName() []byte {
	return d.getKV(true)
}

// getKV gets key-value string for DelayOutput
func (d *DelayOutput) getKV(withProgramName bool) []byte {
	// convert to key-value string
	var buf bytes.Buffer

	// current
	str := d.keyPrefixGen("Current", withProgramName)
	d.Current.KVString(&buf, str)

	// past
	str = d.keyPrefixGen("Past", withProgramName)
	d.Past.KVString(&buf, str)

	return buf.Bytes()
}

func (d *DelayOutput) GetPrometheusFormat() []byte {
	var buf bytes.Buffer
	str := d.keyPrefixGen("Past", true)
	d.Past.PrometheusString(&buf, str)
	return buf.Bytes()
}
