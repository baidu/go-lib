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
Usage 1: Define Metric Field at first
    import "github.com/baidu/go-lib/web-monitor/metrics"

	// define counter struct type
	type ServerState {
		ReqServed *Counter // field type must be *Counter or *Gauge or *State
		ConServed *Counter
		ConActive *Gauge
	}

	// create metrics
	var m Metrics
    var s ServerState
    m.Init(&s, "PROXY", 20)

	// counter operations
	s.ConActive.Inc(2)
    s.ConServed.Inc(1)
    s.ReqServed.Inc(1)
    s.ConActive.Dec(1)

	// get absoulute data for all metrics
    stateData := m.GetAll()
	// get diff data for all counters(gauge don't have diff data)
    stateDiff := m.GetDiff()

Usage 2: Dynamic LoadMetrics Object(if not exist, create it)
    import "github.com/baidu/go-lib/web-monitor/metrics"
	m := NewMetrics("test", 1)

	m.LoadCounter("COUNT").Inc(1)
	m.LoadGauge("Gauge").Inc(1)
	m.LoadState("State").Set("State")

	// get absoulute data for all metrics
    stateData := m.GetAll()
	// get diff data for all counters(gauge don't have diff data)
    stateDiff := m.GetDiff()
*/
package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
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

const (
	TypeGauge   = "Gauge"
	TypeCounter = "Counter"
	TypeState   = "State"
)

var (
	errStructPtrType   = errors.New("metrics should be struct pointor")
	errStructFieldType = errors.New("struct field shoule be *Counter or *Gauge or *State")
)

var (
	supportTypes = map[string]bool{TypeGauge: true, TypeCounter: true, TypeState: true}
)

type Metrics struct {
	// constant after initial
	metricStruct interface{} // underlying struct (pointor)
	metricPrefix string      // name prefix for all metrics
	interval     int         // diff interval
	counterMap   map[string]*Counter
	gaugeMap     map[string]*Gauge
	stateMap     map[string]*State

	// protect following fields
	lock        sync.RWMutex
	metricsLast *MetricsData // last absolute counters
	metricsDiff *MetricsData // diff in last duration
}

// NewMetrics new metrics object, you can got one empty Metrics if you want dynamic get/create metrics by name
func NewMetrics(prefix string, intervalS int) *Metrics {
	m := &Metrics{}
	m.Init(&struct{}{}, prefix, intervalS)

	return m
}

// LoadCounter load counter by name, if not existed, new count will be created then return
func (m *Metrics) LoadCounter(name string) *Counter {
	key := m.convert(name)

	m.lock.RLock()
	if val, ok := m.counterMap[key]; ok {
		m.lock.RUnlock()
		return val
	}
	m.lock.RUnlock()

	m.lock.Lock()
	defer m.lock.Unlock()

	if val, ok := m.counterMap[key]; ok {
		return val
	}

	val := new(Counter)
	m.counterMap[key] = val
	return val
}

// LoadGauge load Gauge by name, if not existed, new count will be created then return
func (m *Metrics) LoadGauge(name string) *Gauge {
	key := m.convert(name)

	m.lock.RLock()
	if val, ok := m.gaugeMap[key]; ok {
		m.lock.RUnlock()
		return val
	}
	m.lock.RUnlock()

	m.lock.Lock()
	defer m.lock.Unlock()

	if val, ok := m.gaugeMap[key]; ok {
		return val
	}

	val := new(Gauge)
	m.gaugeMap[key] = val
	return val
}

// LoadState load state by name, if not existed, new count will be created then return
func (m *Metrics) LoadState(name string) *State {
	key := m.convert(name)

	m.lock.RLock()
	if val, ok := m.stateMap[key]; ok {
		m.lock.RUnlock()
		return val
	}
	m.lock.RUnlock()

	m.lock.Lock()
	defer m.lock.Unlock()

	if val, ok := m.stateMap[key]; ok {
		return val
	}

	val := new(State)
	m.stateMap[key] = val
	return val
}

// Init initializes metrics
//
// Params:
//     - counters: a pointer to a sturct var; struct field type must be *Counter
//     - prefix  : prefix for counters
//     - interval: diff interval (second), if <=0, use default value 20
func (m *Metrics) Init(metrics interface{}, prefix string, interval int) error {
	if err := validateMetrics(metrics); err != nil {
		return err
	}
	if interval <= 0 {
		interval = DefaultInterval
	}

	m.metricStruct = metrics
	m.metricPrefix = prefix
	m.interval = interval
	m.initMetrics(metrics)

	// zero all counters
	m.metricsLast = m.GetAll()
	m.metricsDiff = m.metricsLast.Diff(m.metricsLast)

	go m.handleCounterDiff(m.interval)
	return nil
}

func validateMetrics(metrics interface{}) error {
	// check type of counters is pointer to struct
	t := reflect.TypeOf(metrics)
	if t.Kind() != reflect.Ptr {
		return errStructPtrType
	}

	s := t.Elem()
	if s.Kind() != reflect.Struct {
		return errStructPtrType
	}

	// check type of struct field is *Counter || *Gauge
	for i := 0; i < s.NumField(); i++ {
		ft := s.Field(i).Type
		if ft.Kind() != reflect.Ptr {
			return errStructFieldType
		}

		fn := ft.Elem().Name()
		if _, ok := supportTypes[fn]; !ok {
			return errStructFieldType
		}
	}

	return nil
}

// GetAll gets absoulute values for all counters
func (m *Metrics) GetAll() *MetricsData {
	d := NewMetricsData(m.metricPrefix, KindTotal)
	for k, c := range m.counterMap {
		d.CounterData[k] = int64(c.Get())
	}

	for k, c := range m.gaugeMap {
		d.GaugeData[k] = int64(c.Get())
	}

	for k, s := range m.stateMap {
		d.StateData[k] = s.Get()
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

// initMetrics initializes metrics struct
func (m *Metrics) initMetrics(s interface{}) {
	m.counterMap = make(map[string]*Counter)
	m.gaugeMap = make(map[string]*Gauge)
	m.stateMap = make(map[string]*State)

	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// track created counters
		name := m.convert(field.Name)

		switch mType := field.Type.Elem().Name(); mType {
		case TypeState:
			v := new(State)
			m.stateMap[name] = v
			value.Set(reflect.ValueOf(v))

		case TypeCounter:
			v := new(Counter)
			m.counterMap[name] = v
			value.Set(reflect.ValueOf(v))

		case TypeGauge:
			v := new(Gauge)
			m.gaugeMap[name] = v
			value.Set(reflect.ValueOf(v))
		}
	}
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
	Prefix      string
	Kind        string
	GaugeData   map[string]int64
	CounterData map[string]int64
	StateData   map[string]string
}

func NewMetricsData(prefix string, kind string) *MetricsData {
	d := new(MetricsData)
	d.Prefix = prefix
	d.Kind = kind
	d.GaugeData = make(map[string]int64)
	d.CounterData = make(map[string]int64)
	d.StateData = make(map[string]string)
	return d
}

func (d *MetricsData) Diff(last *MetricsData) *MetricsData {
	diff := NewMetricsData(d.Prefix+"_diff", KindDelta)

	for k, v := range d.CounterData {

		if v2, ok := last.CounterData[k]; ok {
			diff.CounterData[k] = v - v2
		} else {
			diff.CounterData[k] = v
		}

	}
	return diff
}

func (d *MetricsData) Sum(d2 *MetricsData) *MetricsData {
	for k, v := range d2.CounterData {
		if v0, ok := d.CounterData[k]; ok {
			d.CounterData[k] = v0 + v
		} else {
			d.CounterData[k] = v
		}
	}
	return d
}

func (d *MetricsData) KeyValueFormat() []byte {
	var b bytes.Buffer
	for k, v := range d.CounterData {
		line := fmt.Sprintf("%s_%s: %d\n", d.Prefix, k, v)
		b.WriteString(line)
	}

	for k, v := range d.GaugeData {
		line := fmt.Sprintf("%s_%s: %d\n", d.Prefix, k, v)
		b.WriteString(line)
	}

	for k, v := range d.StateData {
		line := fmt.Sprintf("%s_%s: %s\n", d.Prefix, k, v)
		b.WriteString(line)
	}
	return b.Bytes()
}

func (d *MetricsData) PrometheusFormat() []byte {
	var b bytes.Buffer
	for k, v := range d.CounterData {
		key := fmt.Sprintf("%s_%s", d.Prefix, k)
		line := fmt.Sprintf("# TYPE %s %s\n%s %d\n", key, strings.ToLower(TypeCounter), key, v)
		b.WriteString(line)
	}

	for k, v := range d.GaugeData {
		key := fmt.Sprintf("%s_%s", d.Prefix, k)
		line := fmt.Sprintf("# TYPE %s %s\n%s %d\n", key, strings.ToLower(TypeGauge), key, v)
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
	case "kv", "noah":
		return d.KeyValueFormat(), nil
	case "prometheus":
		return d.PrometheusFormat(), nil
	default:
		return nil, fmt.Errorf("invalid format: %s", format)
	}
}
