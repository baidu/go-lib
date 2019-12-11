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

// register web handlers

package web_monitor

import (
	"errors"
	"fmt"
)

// RegisterHandlers registers handlers in handler-table to WebHandlers
//
// Params:
//      - wh    : WebHandlers
//      - hType : hanlder type, WebHandleMonitor or WebHandleReload or WebHandlePprof
//      - ht    : handler table
func RegisterHandlers(wh *WebHandlers, hType int, ht map[string]interface{}) error {
	// check WebHandlers
	if wh == nil {
		return errors.New("nil WebHandlers")
	}

	// check hType
	var typeStr string
	switch hType {
	case WebHandleMonitor:
		typeStr = "MONITOR"
	case WebHandleReload:
		typeStr = "RELOAD"
	case WebHandlePprof:
		typeStr = "PPROF"
	default:
		return fmt.Errorf("invalid handler type:%d", hType)
	}

	// register handlers
	for name, handler := range ht {
		err := wh.RegisterHandler(hType, name, handler)
		if err != nil {
			return fmt.Errorf("register:%s:%s:%s", typeStr, name, err.Error())
		}
	}

	return nil
}
