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

package metrics

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

import (
	"github.com/baidu/go-lib/web-monitor/module_state2"
)

type MockState struct {
	ModuleCounter0 *Counter
	ModuleCounter1 *Counter
	ModuleCounter2 *Counter
	ModuleCounter3 *Counter
	ModuleCounter4 *Counter
	ModuleCounter5 *Counter
	ModuleCounter6 *Counter
	ModuleCounter7 *Counter
	ModuleCounter8 *Gauge
	ModuleCounter9 *State
}

func prepareMetricsState() (*Metrics, *MockState) {
	state := new(MockState)
	metrics := new(Metrics)
	metrics.Init(state, "METRICS", 20)
	return metrics, state
}

func prepareModuleState() module_state2.State {
	counters := []string{
		"MODULE_COUNTER0",
		"MODULE_COUNTER1",
		"MODULE_COUNTER2",
		"MODULE_COUNTER3",
		"MODULE_COUNTER4",
		"MODULE_COUNTER5",
		"MODULE_COUNTER6",
		"MODULE_COUNTER7",
		"MODULE_COUNTER8",
		"MODULE_COUNTER9",
	}
	var s module_state2.State
	s.Init()
	s.CountersInit(counters)
	return s
}

func TestMetricsGetAll(t *testing.T) {
	m, s := prepareMetricsState()
	s.ModuleCounter0.Inc(1)
	s.ModuleCounter2.Inc(1)
	s.ModuleCounter2.Inc(1)
	s.ModuleCounter4.Inc(4)
	s.ModuleCounter8.Dec(2)
	s.ModuleCounter9.Set("version")

	d := m.GetAll()

	r := NewMetricsData("METRICS", KindTotal)
	r.CounterData["MODULE_COUNTER0"] = int64(1)
	r.CounterData["MODULE_COUNTER1"] = int64(0)
	r.CounterData["MODULE_COUNTER2"] = int64(2)
	r.CounterData["MODULE_COUNTER3"] = int64(0)
	r.CounterData["MODULE_COUNTER4"] = int64(4)
	r.CounterData["MODULE_COUNTER5"] = int64(0)
	r.CounterData["MODULE_COUNTER6"] = int64(0)
	r.CounterData["MODULE_COUNTER7"] = int64(0)
	r.GaugeData["MODULE_COUNTER8"] = int64(-2)
	if !reflect.DeepEqual(d.CounterData, r.CounterData) {
		t.Errorf("GetAll(): expect %v, actual %v", r.CounterData, d.CounterData)
	}

	if !reflect.DeepEqual(d.GaugeData, r.GaugeData) {
		t.Errorf("GetAll(): expect %v, actual %v", r.GaugeData, d.GaugeData)
	}

	r.StateData["MODULE_COUNTER9"] = "version"
	if !reflect.DeepEqual(d.StateData, r.StateData) {
		t.Errorf("GetAll():  states expect %v, actual %v", r.StateData, d.StateData)
	}
}

func TestMetricsGetDiff(t *testing.T) {
	m, s := prepareMetricsState()

	// case 1
	d := m.GetDiff()
	r := NewMetricsData("METRICS", KindDelta)
	r.CounterData["MODULE_COUNTER0"] = int64(0)
	r.CounterData["MODULE_COUNTER1"] = int64(0)
	r.CounterData["MODULE_COUNTER2"] = int64(0)
	r.CounterData["MODULE_COUNTER3"] = int64(0)
	r.CounterData["MODULE_COUNTER4"] = int64(0)
	r.CounterData["MODULE_COUNTER5"] = int64(0)
	r.CounterData["MODULE_COUNTER6"] = int64(0)
	r.CounterData["MODULE_COUNTER7"] = int64(0)

	if !reflect.DeepEqual(d.CounterData, r.CounterData) {
		t.Errorf("GetDiff(): expect %v, actual %v", r.CounterData, d.CounterData)
	}

	// case 2
	s.ModuleCounter0.Inc(1)
	s.ModuleCounter4.Inc(4)
	m.updateDiff()
	d = m.GetDiff()
	r = NewMetricsData("METRICS", KindDelta)
	r.CounterData["MODULE_COUNTER0"] = int64(1)
	r.CounterData["MODULE_COUNTER1"] = int64(0)
	r.CounterData["MODULE_COUNTER2"] = int64(0)
	r.CounterData["MODULE_COUNTER3"] = int64(0)
	r.CounterData["MODULE_COUNTER4"] = int64(4)
	r.CounterData["MODULE_COUNTER5"] = int64(0)
	r.CounterData["MODULE_COUNTER6"] = int64(0)
	r.CounterData["MODULE_COUNTER7"] = int64(0)

	if !reflect.DeepEqual(d.CounterData, r.CounterData) {
		t.Errorf("GetAll(): expect %v, actual %v", r, d)
	}
}

type CaseStructA struct {
	c *Counter
}

type CaseStructB struct {
	c Counter
}

type CaseStructC struct {
	c *Counter
	I *CaseStructA
}

func TestInvalidCounter(t *testing.T) {
	var m Metrics

	// case 1
	var s1 CaseStructA
	if err := m.Init(s1, "METRICS", 20); err == nil {
		t.Errorf("expect error: %s", errStructPtrType)
	}

	// case 2
	var s2 CaseStructB
	if err := m.Init(&s2, "METRICS", 20); err == nil {
		t.Errorf("expect error: %s", errStructFieldType)
	}

	// case 3
	var s3 CaseStructC
	if err := m.Init(&s3, "METRICS", 20); err == nil {
		t.Errorf("expect error: %s", errStructFieldType)
	}
}

func TestMetricsData(t *testing.T) {
	d1 := NewMetricsData("METRIX", KindTotal)
	d2 := NewMetricsData("METRIX", KindTotal)
	de := NewMetricsData("METRIX", KindDelta)

	d1.CounterData["MODULE_COUNTER1"] = 10
	d1.CounterData["MODULE_COUNTER2"] = 20
	d2.CounterData["MODULE_COUNTER2"] = 30
	d2.CounterData["MODULE_COUNTER3"] = 40
	de.CounterData["MODULE_COUNTER2"] = 10
	de.CounterData["MODULE_COUNTER3"] = 40

	d := d2.Diff(d1)
	if !reflect.DeepEqual(de.CounterData, d.CounterData) {
		t.Errorf("expect %v, actual %v", de.CounterData, d.CounterData)
	}
}

func BenchmarkMetricIncSingle(b *testing.B) {
	_, s := prepareMetricsState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ModuleCounter4.Inc(1)
	}
}

func BenchmarkStateIncSingle(b *testing.B) {
	s := prepareModuleState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Inc("MODULE_COUNTER4", 1)
	}
}

func BenchmarkMetricIncMulti(b *testing.B) {
	_, s := prepareMetricsState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ModuleCounter4.Inc(1)
		s.ModuleCounter0.Inc(1)
		s.ModuleCounter6.Inc(1)
	}
}

func BenchmarkStateIncMutli(b *testing.B) {
	s := prepareModuleState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Inc("MODULE_COUNTER4", 1)
		s.Inc("MODULE_COUNTER0", 1)
		s.Inc("MODULE_COUNTER6", 1)
	}
}

func BenchmarkMetricGet(b *testing.B) {
	m, _ := prepareMetricsState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.GetAll()
	}
}

func BenchmarkStateGet(b *testing.B) {
	s := prepareModuleState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetAll()
	}
}

func TestNewMetrics(t *testing.T) {
	m := NewMetrics("test", 1)

	tick200Ms := time.NewTicker(time.Millisecond * 200)
	tick2S := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-tick200Ms.C:
			count := len(m.GetAll().CounterData)
			if count !=0 {
				if count != 5 {
					t.Errorf("CountData len want 0/5, got %v", count)
				}else {
					for i := 0; i < 5; i++ {
						key := fmt.Sprintf("COUNT_%d", i)
						if v := m.GetAll().CounterData[key]; v < 0 {
							t.Errorf("CountData[%v] want > 0,  got %v", key, v)
						}
					}
				}

			}

			count = len(m.GetAll().GaugeData)
			if count !=0 && count != 5 {
				t.Errorf("GaugeData len want 0/5, got %v", count)
			}
			count = len(m.GetAll().StateData)
			if count !=0 && count != 5 {
				t.Errorf("StateData len want 0/5, got %v", count)
			}
			// no more strict validate for XxxData[key]

			// simulate concurrent
			go func() {
				for i := 0; i <= 100; i++ {
					m.LoadCounter(fmt.Sprintf("COUNT_%d", i%5)).Inc(1)
					m.LoadGauge(fmt.Sprintf("GAUGE_%d", i%5)).Inc(1)
					m.LoadState(fmt.Sprintf("State_%d", i%5)).Set("State")
				}

				for i := 100; i > 0; i-- {
					m.LoadCounter(fmt.Sprintf("COUNT_%d", i%5)).Inc(1)
					m.LoadGauge(fmt.Sprintf("GAUGE_%d", i%5)).Inc(1)
					m.LoadState(fmt.Sprintf("State_%d", i%5)).Set("State")
				}
			}()
		case <-tick2S.C:
			tick2S.Stop()
			tick200Ms.Stop()
			return
		}
	}

}
