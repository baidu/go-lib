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

func TestCounterGet(t *testing.T) {
	var c Counter
	if c.Get() != 0 {
		t.Errorf("init counter expect 0, but is:%d", c.Get())
	}
}

func TestCounterInc(t *testing.T) {
	var c Counter
	c.Inc(10)
	if c.Get() != 10 {
		t.Errorf("after inc 10, counter expected to be 10, but is:%d", c.Get())
	}
}
