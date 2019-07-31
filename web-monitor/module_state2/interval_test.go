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

import (
	"fmt"
	"testing"
	"time"
)

func TestNextInterval(t *testing.T) {
	// test case 1
	now := time.Date(2009, time.November, 10, 23, 10, 10, 0, time.UTC)
	interval := NextInterval(now, 60)
	if interval != 50 {
		t.Error(fmt.Sprintf("return of NextInterval() should be 50, it's %d", interval))
	}

	// test case 2
	now = time.Date(2009, time.November, 10, 23, 10, 0, 0, time.UTC)
	interval = NextInterval(now, 60)
	if interval != 60 {
		t.Error(fmt.Sprintf("return of NextInterval() should be 60, it's %d", interval))
	}

	// test case 3
	now = time.Date(2009, time.November, 10, 23, 10, 40, 0, time.UTC)
	interval = NextInterval(now, 60)
	if interval != 20 {
		t.Error(fmt.Sprintf("return of NextInterval() should be 20, it's %d", interval))
	}
}
