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
		ReqServed *Counter // field type must be *Counter or *Gauge
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
)

var (
	errStructPtrType   = errors.New("metrics should be struct pointor")
	errStructFieldType = errors.New("struct field shoule be *Counter or *Gauge")
)

var (
	supportTypes = map[string]bool{TypeGauge: true, TypeCounter: true}
)

type Metric interface {
	Get() int64
	Type() string
}

type Metrics struct {
	// constant after initial
	metricStruct interface{}       // underlying struct (pointor)
	metricPrefix string            // name prefix for all metrics
	interval     int               // diff interval
	metricsMap   map[string]Metric // all metrics

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
	for k, c := range m.metricsMap {
		d.DataTypes[k] = c.Type()
		d.Data[k] = int64(c.Get())
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
	m.metricsMap = make(map[string]Metric)
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// track created counters
		name := m.convert(field.Name)

		var metric Metric

		switch mType := field.Type.Elem().Name(); mType {
		case TypeCounter:
			metric = new(Counter)
		case TypeGauge:
			metric = new(Gauge)
		}

		m.metricsMap[name] = metric
		value.Set(reflect.ValueOf(metric))
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
	Prefix    string
	Kind      string
	DataTypes map[string]string
	Data      map[string]int64
}

func NewMetricsData(prefix string, kind string) *MetricsData {
	d := new(MetricsData)
	d.Prefix = prefix
	d.Kind = kind
	d.DataTypes = make(map[string]string)
	d.Data = make(map[string]int64)
	return d
}

func (d *MetricsData) Diff(last *MetricsData) *MetricsData {
	diff := NewMetricsData(d.Prefix+"_diff", KindDelta)
	diff.DataTypes = d.DataTypes

	for k, v := range d.Data {
		dt, _ := d.DataTypes[k]
		if dt == TypeGauge {
			continue
		}

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

func (d *MetricsData) KeyValueFormat() []byte {
	var b bytes.Buffer
	for k, v := range d.Data {
		line := fmt.Sprintf("%s_%s: %d\n", d.Prefix, k, v)
		b.WriteString(line)
	}
	return b.Bytes()
}

func (d *MetricsData) PrometheusFormat() []byte {
	var b bytes.Buffer
	for k, v := range d.Data {
		key := fmt.Sprintf("%s_%s", d.Prefix, k)
		line := fmt.Sprintf("# TYPE %s %s\n%s %d\n", key, strings.ToLower(d.DataTypes[k]), key, v)
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
		return d.KeyValueFormat(), nil
	case "prometheus":
		return d.PrometheusFormat(), nil
	default:
		return nil, fmt.Errorf("invalid format: %s", format)
	}
}
