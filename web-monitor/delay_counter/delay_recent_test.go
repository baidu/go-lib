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

package delay_counter

import (
	"testing"
	"time"
)

import (
	"github.com/baidu/go-lib/log"
)

func TestDelayRecent(t *testing.T) {
	log.Init("test", "DEBUG", "./log", true, "D", 5)

	var delayTable DelayRecent

	// initialize the table
	// interval=20, bucketSize=1, bucketNum=10
	delayTable.Init(20, 1, 10)

	start := time.Now()

	// try to get when table is empty
	_, err1 := delayTable.GetJson()
	if err1 != nil {
		t.Error("Error in DelayTableGet()")
	}

	// try to invoke sub() of DelayTable
	delayTable.AddBySub(start, time.Now())

	// try to get again
	_, err2 := delayTable.GetJson()
	if err2 != nil {
		t.Error("Error in DelayTableGet()")
	}

	// try to invoke add() of DelayTable
	duration := time.Now().Sub(start).Nanoseconds() / 1000
	delayTable.Add(duration)

	// try to get again
	_, err2 = delayTable.GetJson()
	if err2 != nil {
		t.Error("Error in DelayTableGet()")
	}

	log.Logger.Close()
}

func TestFormatOutput(t *testing.T) {
	var delay DelayRecent
	delay.Init(20, 1, 100)

	params := map[string][]string{
		"format": []string{"json"},
	}
	_, err := (&delay).FormatOutput(params)
	if err != nil {
		t.Errorf("FormatOutDR(): testcase 0 : %s", err.Error())
	}

	params = map[string][]string{
		"format": []string{"kv"},
	}

	_, err = (&delay).FormatOutput(params)
	if err != nil {
		t.Errorf("FormatOutDR(): testcase 1 : %s", err.Error())
	}

	params = map[string][]string{
		"format": []string{"no_kv"},
	}

	_, err = (&delay).FormatOutput(params)
	if err == nil {
		t.Errorf("FormatOutDR(): testcase 2 should return error!")
	}
}
