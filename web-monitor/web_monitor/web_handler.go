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

// web handler framework

package web_monitor

import (
	"fmt"
	"net/url"
)

// type of web handler
const (
	WEB_HANDLE_MONITOR = 0 // handler for monitor
	WEB_HANDLE_RELOAD  = 1 // handler for reload
)

var handlerTypeNames = map[int]string{
	0: "monitor",
	1: "reload",
}

type WebHandlerMap map[string]interface{}

type WebHandlers struct {
	Handlers map[int]*WebHandlerMap
}

// NewWebHandlerMap creates new WebHandlerMap
func NewWebHandlerMap() *WebHandlerMap {
	whm := make(WebHandlerMap)
	return &whm
}

// NewWebHandlers creates new WebHandlers
func NewWebHandlers() *WebHandlers {
	// create bfeCallbacks
	wh := new(WebHandlers)
	wh.Handlers = make(map[int]*WebHandlerMap)

	// handlers for monitor
	wh.Handlers[WEB_HANDLE_MONITOR] = NewWebHandlerMap()
	// handlers for reload
	wh.Handlers[WEB_HANDLE_RELOAD] = NewWebHandlerMap()

	return wh
}

func (wh *WebHandlers) validateHandler(hType int, f interface{}) error {
	var err error
	switch hType {
	case WEB_HANDLE_MONITOR:
		switch f.(type) {
		case func() ([]byte, error):
		case func(map[string][]string) ([]byte, error):
		case func(url.Values) ([]byte, error):
		default:
			err = fmt.Errorf("invalid monitor handler type %T", f)
		}

	case WEB_HANDLE_RELOAD:
		switch f.(type) {
		case func() error:
		case func(map[string][]string) error:
		case func(url.Values) error:
		case func(url.Values) (string, error):
		default:
			err = fmt.Errorf("invalid reload handler type %T", f)
		}

	default:
		err = fmt.Errorf("invalid handler type[%d]", hType)
	}
	return err
}

// RegisterHandler adds filter to given callback point
func (wh *WebHandlers) RegisterHandler(hType int, command string, f interface{}) error {
	var ok bool
	var hm *WebHandlerMap

	// check format of f
	if err := wh.validateHandler(hType, f); err != nil {
		return err
	}

	// get WebHandlerMap for given hType
	hm, ok = wh.Handlers[hType]
	if !ok {
		return fmt.Errorf("invalid handler type[%d]", hType)
	}

	// handler exist already?
	_, ok = (*hm)[command]
	if ok {
		return fmt.Errorf("handler exist already, type[%s], command[%s]",
			handlerTypeNames[hType], command)
	}

	// add to WebHandlerMap
	(*hm)[command] = f

	return nil
}

// GetHandler gets handler list for given callback point
func (wh *WebHandlers) GetHandler(hType int, command string) (interface{}, error) {
	var ok bool
	var hm *WebHandlerMap
	var h interface{}

	// get WebHandlerMap for given hType
	hm, ok = wh.Handlers[hType]
	if !ok {
		return nil, fmt.Errorf("invalid handler type[%d]", hType)
	}

	// handler exist already?
	h, ok = (*hm)[command]
	if !ok {
		return nil, fmt.Errorf("handler not exist, type[%s], command[%s]",
			handlerTypeNames[hType], command)
	}

	return h, nil
}
