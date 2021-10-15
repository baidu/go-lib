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
	"testing"
)

func TestGaugeGet(t *testing.T) {
	var g Gauge
	if g.Get() != 0 {
		t.Errorf("init gauge expected 0, but is:%d", g.Get())
	}
}

func TestGaugeInc(t *testing.T) {
	var g Gauge
	g.Inc(10)
	if g.Get() != 10 {
		t.Errorf("after inc 10, gauge expected 10, but is:%d", g.Get())
	}
}

func TestGaugeDec(t *testing.T) {
	var g Gauge
	g.Dec(5)
	if g.Get() != -5 {
		t.Errorf("after dec 5, gauge expected -5, but is:%d", g.Get())
	}
}

func TestGaugeSet(t *testing.T) {
	var g Gauge
	g.Set(3)
	if g.Get() != 3 {
		t.Errorf("after set 3, gauge expected 3, but is:%d", g.Get())
	}
}
