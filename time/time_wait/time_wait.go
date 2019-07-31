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
	"time"
)

// WaitTill waits until toTime
//
// Params:
//     - toTime: time to wait until. the number of seconds elapsed since January 1, 1970 UTC.
func WaitTill(toTime int64) {
	waitSecs := toTime - time.Now().Unix()
	if waitSecs > 0 {
		time.Sleep(time.Second * time.Duration(waitSecs))
	}
}

// CalcNextTime calculates the nearest time from now, given cycle and offset
// 
// Params:
//  - cycle: cycle in seconds
//  - offset: offset of the next time; in seconds
// 
// Return:
//  - timestamp of next time
func CalcNextTime(cycle int64, offset int64) int64 {
	current := time.Now().Unix()

	if current%cycle == 0 {
		return current + offset
	} else {
		return current - current%cycle + cycle + offset
	}
}
