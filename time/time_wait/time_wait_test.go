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

package time_wait

import (
	"testing"
	"time"
)

func TestWaitTill_case1(t *testing.T) {
	waitSecs := int64(2)
	start := time.Now().Unix()

	toTime := start + waitSecs

	WaitTill(toTime)

	passSecs := time.Now().Unix() - start

	if passSecs != waitSecs {
		t.Errorf("err in WaitTill(): wait=%d, pass=%d", waitSecs, passSecs)
	}
}

func TestWaitTill_case2(t *testing.T) {
	start := time.Now().Unix()

	toTime := start - 2

	WaitTill(toTime)

	passSecs := time.Now().Unix() - start

	if passSecs != 0 {
		t.Errorf("err in WaitTill(): wait=0, pass=%d", passSecs)
	}
}
