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

/*
embeded web server for monitor and reload

This web server is for:
- monitor internal state of daemon server
- reload config for daemon server
*/

package web_monitor

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

import (
	"github.com/baidu/go-lib/gotrack"
	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/time/timefmt"
)

// source ip address allowed to do reload
var RELOAD_SRC_ALLOWED = map[string]bool{
	"127.0.0.1": true,
	"::1":       true,
}

type MonitorServer struct {
	port        int          // port for listen
	name        string       // name of the daemon server
	version     string       // version of daemon server
	startAt     string       // start time of daemon server
	webHandlers *WebHandlers // table of web handlers
}

// NewMonitorServer creates new MonitorServer
func NewMonitorServer(name string, version string, port int) *MonitorServer {
	srv := new(MonitorServer)

	srv.name = name
	srv.version = version
	srv.startAt = timefmt.CurrTimeGet()
	srv.port = port

	srv.webHandlers = NewWebHandlers()

	return srv
}

// RegisterHandler registers handler
func (srv *MonitorServer) RegisterHandler(hType int, command string, f interface{}) error {
	var err error

	switch hType {
	case WebHandleMonitor, WebHandleReload, WebHandlePprof:
		err = srv.webHandlers.RegisterHandler(hType, command, f)
	default:
		err = fmt.Errorf("invalid handler type[%d]", hType)
	}

	return err
}

// RegisterHandlers registers handlers in handler-table to WebHandlers
//
// Params:
//      - hType : hanlder type, WebHandleMonitor or WebHandleReload or WebHandlePprof
//      - ht    : handler table
func (srv *MonitorServer) RegisterHandlers(hType int, ht map[string]interface{}) error {
	var err error

	switch hType {
	case WebHandleMonitor, WebHandleReload, WebHandlePprof:
		err = RegisterHandlers(srv.webHandlers, hType, ht)
	default:
		err = fmt.Errorf("invalid handler type[%d]", hType)
	}

	return err

}

// HandlersSet sets handlers
func (srv *MonitorServer) HandlersSet(handlers *WebHandlers) {
	srv.webHandlers = handlers
}

func abnormalExit() {
	/* to overcome bug in log, sleep for a while    */
	time.Sleep(1 * time.Second)
	os.Exit(1)
}

// isValidForReload checks whether remote address is valid for doing reload
func isValidForReload(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	_, ok := RELOAD_SRC_ALLOWED[host]
	return ok
}

// subManualShow shows sub manual (monitor/reload)
func (srv *MonitorServer) subManualShow(hType int) []byte {
	var commands = make([]string, 0)

	for k, _ := range *srv.webHandlers.Handlers[hType] {
		commands = append(commands, k)
	}
	sort.Strings(commands)

	var typeStr string
	switch hType {
	case WebHandleMonitor:
		typeStr = "monitor"
	case WebHandleReload:
		typeStr = "reload"
	case WebHandlePprof:
		typeStr = "debug"
	}

	str := "<html>\n"
	str += "<body>\n"
	str += fmt.Sprintf("<p>%s manual for %s</p>\n", typeStr, srv.name)
	str += fmt.Sprintf("<p>version: %s</p>\n", srv.version)
	str += fmt.Sprintf("<p>start_at: %s</p>\n", srv.startAt)

	for _, command := range commands {
		line := fmt.Sprintf("<p><a href=\"/%s/%s\">%s</a></p>\n", typeStr, command, command)
		str = str + line
	}

	str += "</body>"
	str += "</html>"

	return []byte(str)
}

// manualShow shows manual of web server
func (srv *MonitorServer) manualShow() []byte {
	str := "<html>\n"
	str += "<body>\n"
	str += fmt.Sprintf("<p>Welcome to %s</p>\n", srv.name)
	str += fmt.Sprintf("<p>version: %s</p>\n", srv.version)
	str += fmt.Sprintf("<p>start_at: %s</p>\n", srv.startAt)
	str = str + fmt.Sprintf("<p><a href=\"/monitor\">monitor</a></p>\n")
	str = str + fmt.Sprintf("<p><a href=\"/reload\">reload</a></p>\n")
	str = str + fmt.Sprintf("<p><a href=\"/debug\">debug</a></p>\n")

	str += "</body>"
	str += "</html>"

	return []byte(str)
}

func errInfoGen(err error) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
}

func webOutput(w http.ResponseWriter, buff []byte, err error) {
	if err != nil {
		errStr := errInfoGen(err)
		fmt.Fprintf(w, "%s", errStr)
	} else {
		fmt.Fprintf(w, "%s", buff)
	}
}

func (srv *MonitorServer) monitorHandler(command string,
	params map[string][]string) (buff []byte, err error) {
	var f interface{}

	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("monitor panic:%v", perr)
			log.Logger.Warn("MonitorServer:monitorHandler():%v\n%s",
				perr, gotrack.CurrentStackTrace(0))
		}
	}()

	// get handler
	f, err = srv.webHandlers.GetHandler(WebHandleMonitor, command)
	if err != nil {
		return buff, err
	}

	// invoke handler for monitor
	switch f.(type) {
	case func() ([]byte, error):
		buff, err = f.(func() ([]byte, error))()
	case func(map[string][]string) ([]byte, error):
		buff, err = f.(func(map[string][]string) ([]byte, error))(params)
	case func(url.Values) ([]byte, error):
		buff, err = f.(func(url.Values) ([]byte, error))(params)
	}

	return buff, err
}

func (srv *MonitorServer) reloadHandler(command string, params map[string][]string,
	remoteAddr string) (buff []byte, err error) {
	var f interface{}
	var version string

	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("reload panic:%v", perr)
			log.Logger.Warn("MonitorServer:reloadHandler():%v\n%s",
				perr, gotrack.CurrentStackTrace(0))
		}
	}()

	// check source address
	if !isValidForReload(remoteAddr) {
		err = fmt.Errorf("reload is not allowed from [%s]", remoteAddr)
		log.Logger.Warn("MonitorServer:Blocked reload request from[%s], cmd=[%s]",
			remoteAddr, command)
		return buff, err
	}

	// get handler
	f, err = srv.webHandlers.GetHandler(WebHandleReload, command)
	if err != nil {
		return buff, err
	}

	// invoke handler for reload
	switch f.(type) {
	case func() error:
		err = f.(func() error)()
	case func(map[string][]string) error:
		err = f.(func(map[string][]string) error)(params)
	case func(url.Values) error:
		err = f.(func(url.Values) error)(params)
	case func(url.Values) (string, error):
		// format of returned version info is like f1=v1&f2=v2, e.g.,
		// host_rule.data=201708280900&route_rule.data=201708280900
		version, err = f.(func(url.Values) (string, error))(params)
	}

	if err != nil {
		log.Logger.Error("MonitorServer:Reload through web, "+
			"cmd=[%s], params=[%s], from[%s], err=[%s]",
			command, remoteAddr, params, err.Error())
		return buff, err
	}

	log.Logger.Info("MonitorServer:Reload through web, cmd=[%s], params=[%s] from[%s]",
		command, params, remoteAddr)

	if version != "" {
		buff = []byte(fmt.Sprintf("{\"error\":null,\"version\":%q}", version))
	} else {
		buff = []byte(fmt.Sprintf("{\"error\":null}"))
	}

	return buff, nil
}

func (srv *MonitorServer) pprofHandler(command string, w http.ResponseWriter, r *http.Request) (err error) {
	var f interface{}

	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("monitor panic:%v", perr)
			log.Logger.Warn("MonitorServer:pprofHandler():%v\n%s",
				perr, gotrack.CurrentStackTrace(0))
		}
	}()

	// get handler
	f, err = srv.webHandlers.GetHandler(WebHandlePprof, command)
	if err != nil {
		return err
	}

	// invoke handler for monitor
	switch f.(type) {
	case func(w http.ResponseWriter, r *http.Request):
		f.(func(w http.ResponseWriter, r *http.Request))(w, r)
	}

	return err
}

func (srv *MonitorServer) webHandler(w http.ResponseWriter, r *http.Request) {
	var buff []byte
	var err error
	var commands []string

	// Path should be:
	//     /monitor/host_table
	//     /reload/mod_trust_clientip
	//
	command := r.URL.Path[1:]
	if len(command) == 0 {
		commands = make([]string, 0)
	} else {
		commands = strings.SplitN(command, "/", 2)
	}
	params := r.URL.Query()

	switch len(commands) {
	case 1:
		switch commands[0] {
		case "monitor":
			buff = srv.subManualShow(WebHandleMonitor)
			err = nil
		case "reload":
			buff = srv.subManualShow(WebHandleReload)
			err = nil
		case "debug":
			buff = srv.subManualShow(WebHandlePprof)
			err = nil
		default:
			err = fmt.Errorf("invalid command [%s]", commands[0])
		}
	case 2:
		switch commands[0] {
		case "monitor":
			buff, err = srv.monitorHandler(commands[1], params)
		case "reload":
			buff, err = srv.reloadHandler(commands[1], params, r.RemoteAddr)
		case "debug":
			err = srv.pprofHandler(commands[1], w, r)
		default:
			err = fmt.Errorf("invalid command [%s]", commands[0])
		}
	default:
		// format error, show the manual
		buff = srv.manualShow()
		err = nil
	}
	webOutput(w, buff, err)
}

// Start starts embeded web server
func (srv *MonitorServer) Start() {
	err := srv.ListenAndServe()
	if err != nil {
		log.Logger.Error("MonitorServer.Start():err in http.ListenAndServe():%s", err.Error())
		abnormalExit()
	}
}

// ListenAndServe start embeded web server
func (srv *MonitorServer) ListenAndServe() error {
	log.Logger.Info("Embeded web server start at port[%d]", srv.port)

	http.HandleFunc("/", srv.webHandler)

	portStr := fmt.Sprintf(":%d", srv.port)
	return http.ListenAndServe(portStr, nil)
}
