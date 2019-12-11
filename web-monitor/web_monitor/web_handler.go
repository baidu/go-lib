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
	"net/http"
	"net/http/pprof"
	"net/url"
)

// type of web handler
const (
	WebHandleMonitor = 0 // handler for monitor
	WebHandleReload  = 1 // handler for reload
	WebHandlePprof   = 2 // handler for pprof
)

var handlerTypeNames = map[int]string{
	0: "monitor",
	1: "reload",
	2: "debug",
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

func pprofHandlers() *WebHandlerMap {
	handlers := &WebHandlerMap{
		"pprof":   pprof.Index,
		"cmdline": pprof.Cmdline,
		"profile": pprof.Profile,
		"symbol":  pprof.Symbol,
		"trace":   pprof.Trace,
	}
	return handlers
}

// NewWebHandlers creates new WebHandlers
func NewWebHandlers() *WebHandlers {
	// create bfeCallbacks
	wh := new(WebHandlers)
	wh.Handlers = make(map[int]*WebHandlerMap)

	// handlers for monitor
	wh.Handlers[WebHandleMonitor] = NewWebHandlerMap()
	// handlers for reload
	wh.Handlers[WebHandleReload] = NewWebHandlerMap()
	// handlers for pprof
	wh.Handlers[WebHandlePprof] = pprofHandlers()

	return wh
}

func (wh *WebHandlers) validateHandler(hType int, f interface{}) error {
	var err error
	switch hType {
	case WebHandleMonitor:
		switch f.(type) {
		case func() ([]byte, error):
		case func(map[string][]string) ([]byte, error):
		case func(url.Values) ([]byte, error):
		default:
			err = fmt.Errorf("invalid monitor handler type %T", f)
		}

	case WebHandleReload:
		switch f.(type) {
		case func() error:
		case func(map[string][]string) error:
		case func(url.Values) error:
		case func(url.Values) (string, error):
		default:
			err = fmt.Errorf("invalid reload handler type %T", f)
		}
	case WebHandlePprof:
		switch f.(type) {
		case func(w http.ResponseWriter, r *http.Request):
		default:
			err = fmt.Errorf("invalid pprof handler type %T", f)
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
