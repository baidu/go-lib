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
	"sync/atomic"
)

// Gauge is a Metric that represents a single numerical value that can
// arbitrarily go up and down.
type Gauge int64

// Inc increases gauge
func (c *Gauge) Inc(delta uint) {
	if c == nil {
		return
	}
	atomic.AddInt64((*int64)(c), int64(delta))
}

// Dec decreases gauge
func (c *Gauge) Dec(delta uint) {
	if c == nil {
		return
	}
	atomic.AddInt64((*int64)(c), int64(-delta))
}

// Get gets gauge
func (c *Gauge) Get() int64 {
	if c == nil {
		return 0
	}
	return atomic.LoadInt64((*int64)(c))
}

// Set sets gauge
func (c *Gauge) Set(v int64) {
	if c == nil {
		return
	}
	atomic.StoreInt64((*int64)(c), v)
}

func (c *Gauge) Type() string {
	return TypeGauge
}
