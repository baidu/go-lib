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

package module_state2

// Counters holds counters for given key
type Counters map[string]int64

// NewCounters creates new Counters
func NewCounters() Counters {
	counters := make(Counters)
	return counters
}

// inc increases value for given key
func (c *Counters) inc(key string, value int) {
	_, ok := (*c)[key]

	if !ok {
		(*c)[key] = int64(value)
	} else {
		(*c)[key] += int64(value)
	}
}

// dec decreases value for given key
func (c *Counters) dec(key string, value int) {
	_, ok := (*c)[key]

	if !ok {
		(*c)[key] = int64(value) * -1
	} else {
		(*c)[key] -= int64(value)
	}
}

// init initializes counter for given keys to zero
func (c *Counters) init(keys []string) {
	for _, key := range keys {
		(*c)[key] = 0
	}
}

// copy makes a copy of Counters
func (c *Counters) copy() Counters {
	copy := make(Counters)
	for key, value := range *c {
		copy[key] = value
	}
	return copy
}

// diff gets change between two counters
func (c *Counters) diff(last Counters) Counters {
	diff := make(Counters)

	for key, value := range *c {
		old_value, ok := last[key]
		if ok {
			// exist in last
			diff[key] = value - old_value
		} else {
			diff[key] = value
		}
	}
	return diff
}

// Sum calculates sum of two Counters
func (c *Counters) Sum(c2 Counters) {
	for key, value2 := range c2 {
		value, ok := (*c)[key]
		if ok {
			(*c)[key] = value + value2
		} else {
			(*c)[key] = value2
		}
	}
}
