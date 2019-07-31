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

/*
Usage:
    import "github.com/baidu/go-lib/web-monitor/metrics"

	// define counter struct type
	type ServerState {
		ReqServed *Counter // field type must be *Counter
		ConServed *Counter
		ConActive *Counter
	}

	// create metrics
	var m Metrics
    var s ServerState
    m.Init(&s, "PROXY", 20)

	// counter operations
    s.ConServed.Inc(1)
    s.ReqServed.Inc(1)
    s.ConActive.Inc(-1)

	// get absoulute and diff data for all counters
    stateData := m.GetAll()
    stateDiff := m.GetDiff()
*/
package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
	"unicode"
)

import (
	"github.com/baidu/go-lib/web-monitor/web_params"
)

const (
	DefaultInterval = 20 // in seconds
)

const (
	KindTotal = "total"
	KindDelta = "delta"
)

var (
	ErrStructPtrType   = errors.New("counters should be struct pointor")
	ErrStructFieldType = errors.New("struct field shoule be *Counter")
)

type Metrics struct {
	// constant after initial
	countersStruct interface{}         // underlying counters struct (pointor)
	countersPrefix string              // name prefix for all conters
	interval       int                 // diff interval
	countersMap    map[string]*Counter // all counters

	// protect following fields
	lock        sync.RWMutex
	metricsLast *MetricsData // last absolute counters
	metricsDiff *MetricsData // diff in last duration
}

// Init initializes metrics
//
// Params:
//     - counters: a pointer to a sturct var; struct field type must be *Counter
//     - prefix  : prefix for counters
//     - interval: diff interval (second), if <=0, use default value 20
func (m *Metrics) Init(counters interface{}, prefix string, interval int) error {
	if err := validateCounters(counters); err != nil {
		return err
	}
	if interval <= 0 {
		interval = DefaultInterval
	}

	m.countersStruct = counters
	m.countersPrefix = prefix
	m.interval = interval
	m.countersMap = m.initCounters(counters)

	// zero all counters
	m.metricsLast = m.GetAll()
	m.metricsDiff = m.metricsLast.Diff(m.metricsLast)

	go m.handleCounterDiff(m.interval)
	return nil
}

func validateCounters(counters interface{}) error {
	// check type of counters is pointer to struct
	t := reflect.TypeOf(counters)
	if t.Kind() != reflect.Ptr {
		return ErrStructPtrType
	}
	s := t.Elem()
	if s.Kind() != reflect.Struct {
		return ErrStructPtrType
	}

	// check type of struct field is *Counter
	var c *Counter
	k := reflect.TypeOf(c).Kind()

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)

		if field.Type.Kind() != k {
			return ErrStructFieldType
		}
	}

	return nil
}

// GetAll gets absoulute values for all counters
func (m *Metrics) GetAll() *MetricsData {
	d := NewMetricsData(m.countersPrefix, KindTotal)
	for k, c := range m.countersMap {
		d.Data[k] = c.Get()
	}
	return d
}

// GetDiff gets diff values for all counters
func (m *Metrics) GetDiff() *MetricsData {
	m.lock.RLock()
	diff := m.metricsDiff
	m.lock.RUnlock()

	return diff
}

// initCounters initializes counters struct
func (m *Metrics) initCounters(s interface{}) map[string]*Counter {
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()
	counters := make(map[string]*Counter)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// track created counters
		name := m.convert(field.Name)
		cntr := new(Counter)
		counters[name] = cntr

		// init struct field
		value.Set(reflect.ValueOf(cntr))
	}

	return counters
}

// convert converts name from CamelCase to UnderScoreCase
func (m *Metrics) convert(name string) string {
	var b bytes.Buffer
	for i, c := range name {
		if unicode.IsUpper(c) {
			if i > 0 {
				b.WriteString("_")
			}
			b.WriteRune(c)
		} else {
			b.WriteRune(unicode.ToUpper(c))
		}
	}
	return b.String()
}

// handleCounterDiff is go-routine for periodically update counter diff
func (m *Metrics) handleCounterDiff(interval int) {
	for {
		m.updateDiff()

		seconds := time.Now().Second()
		left := m.interval - seconds%m.interval
		time.Sleep(time.Duration(left) * time.Second)
	}
}

// updateDiff updates diff values for all counters
func (m *Metrics) updateDiff() {
	var diff *MetricsData

	m.lock.RLock()
	last := m.metricsLast
	m.lock.RUnlock()

	// calc diff data
	current := m.GetAll()
	diff = current.Diff(last)

	// update last and diff data
	m.lock.Lock()
	m.metricsLast = current
	m.metricsDiff = diff
	m.lock.Unlock()
}

type MetricsData struct {
	Prefix string
	Kind   string
	Data   map[string]int64
}

func NewMetricsData(prefix string, kind string) *MetricsData {
	d := new(MetricsData)
	d.Prefix = prefix
	d.Kind = kind
	d.Data = make(map[string]int64)
	return d
}

func (d *MetricsData) Diff(last *MetricsData) *MetricsData {
	diff := NewMetricsData(d.Prefix, KindDelta)

	for k, v := range d.Data {
		if v2, ok := last.Data[k]; ok {
			diff.Data[k] = v - v2
		} else {
			diff.Data[k] = v
		}
	}
	return diff
}

func (d *MetricsData) Sum(d2 *MetricsData) *MetricsData {
	for k, v := range d2.Data {
		if v0, ok := d.Data[k]; ok {
			d.Data[k] = v0 + v
		} else {
			d.Data[k] = v
		}
	}
	return d
}

func (d *MetricsData) Value() []byte {
	p := d.Prefix
	if d.Kind == KindDelta {
		p = p + "_diff"
	}

	var b bytes.Buffer
	for k, v := range d.Data {
		line := fmt.Sprintf("%s_%s:%d\n", p, k, v)
		b.WriteString(line)
	}
	return b.Bytes()
}

func (d *MetricsData) Format(params map[string][]string) ([]byte, error) {
	format, err := web_params.ParamsValueGet(params, "format")
	if err != nil {
		format = "json"
	}

	switch format {
	case "json":
		return json.Marshal(d)
	case "kv":
		return d.Value(), nil
	default:
		return nil, fmt.Errorf("invalid format: %s", format)
	}
}
