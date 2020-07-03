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

package log4go

import (
	"fmt"
)

import (
	"github.com/RackSec/srslog"
)

var (
	syslogger *srslog.Writer
)

func SetSysLogger(network string, addr string, tag string) error {
	w, err := srslog.Dial(network, addr, srslog.LOG_INFO, tag)
	if err != nil {
		return fmt.Errorf("set syslog faile:%v", err)
	}
	syslogger = w
	return nil
}

func ToSyslogPriority(lv LevelType) srslog.Priority {
	switch lv {
	case INFO:
		return srslog.LOG_INFO
	case WARNING:
		return srslog.LOG_WARNING
	case ERROR:
		return srslog.LOG_ERR
	case CRITICAL:
		return srslog.LOG_CRIT
	default:
		return srslog.LOG_DEBUG
	}

}
