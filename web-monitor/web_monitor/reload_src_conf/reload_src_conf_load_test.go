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

package reload_src_conf

import (
	"reflect"
	"sort"
	"testing"
)

func TestReloadSrcIPsLoad(t *testing.T) {
	reloadSrcIPs, err := ReloadSrcIPsLoad("./testdata/reload_src_conf_1.data")
	if err != nil {
		t.Errorf("get err from ReloadSrcIPsLoad():%s", err.Error())
		return
	}

	ips := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}
	sort.Strings(reloadSrcIPs)

	if !reflect.DeepEqual(reloadSrcIPs, ips) {
		t.Errorf("ReloadSrcIPsLoad failed, should be:%v, but is:%v", ips, reloadSrcIPs)
	}
}

func TestReloadSrcIPsLoad_InvalidIP(t *testing.T) {
	_, err := ReloadSrcIPsLoad("./testdata/reload_src_conf_2.data")
	if err == nil {
		t.Fatalf("Expect an error")
	}
}
