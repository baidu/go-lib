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

// conf file parser for reload iplist

package reload_src_conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

type IPList []string // list of ips

type Label2IP map[string]IPList // label => ip list

type ReloadSrcConf struct {
	Version string   // version of the config
	Config  Label2IP // label => ip list
}

// LoadAndCheck loads iplist from filename
func (conf *ReloadSrcConf) LoadAndCheck(filename string) (string, error) {
	// open the file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("os.Open() err:%s", err.Error())
	}
	defer file.Close()

	// decode the file
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(conf); err != nil {
		return "", fmt.Errorf("decoder.Decode() err:%s", err.Error())
	}

	// check config
	if err := ReloadSrcConfCheck(*conf); err != nil {
		return "", fmt.Errorf("ReloadSrcConfCheck() err:%s", err.Error())
	}

	return conf.Version, nil
}

// ReloadSrcConfCheck checks reload conf
func ReloadSrcConfCheck(conf ReloadSrcConf) error {
	if conf.Version == "" {
		return errors.New("no Version")
	}

	// check config for each label
	for label, ipList := range conf.Config {
		var formattedIPList IPList
		for _, ip := range ipList {
			ip2, err := net.ResolveIPAddr("ip", ip)
			if err != nil {
				return fmt.Errorf("invalid ip:%s, in label:%s", ip, label)
			}

			formattedIPList = append(formattedIPList, ip2.String())
		}
		conf.Config[label] = formattedIPList
	}
	return nil
}

// ReloadSrcIPsLoad loads reload iplist allowed from file
func ReloadSrcIPsLoad(filename string) ([]string, error) {
	// load reload src config
	var config ReloadSrcConf
	if _, err := config.LoadAndCheck(filename); err != nil {
		return nil, err
	}

	reloadSrcIPs := make([]string, 0)
	// convert from ReloadSrcConf
	for _, ipList := range config.Config {
		reloadSrcIPs = append(reloadSrcIPs, ipList...)
	}

	return reloadSrcIPs, nil
}
