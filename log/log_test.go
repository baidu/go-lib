// Copyright (c) 2019 Baidu, Inc.
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

package log

import (
    "testing"
    "time"
)

func TestLog(t *testing.T) {
    if err := Init("test", "INFO", "./log/log", true, "M", 2); err != nil {
        t.Error("log.Init() fail")
    }
    
    if err := Init("test", "INFO", "./log/log", true, "M", 5); err == nil {
        t.Error("fail in process reentering log.Init()")
    }

    for i:=0; i < 100; i = i + 1 {
        Logger.Warn("warning msg: %d", i)
        Logger.Info("info msg: %d", i)
        
        // time.Sleep(10 * time.Second)
    }
    
    time.Sleep(100 * time.Millisecond)
}