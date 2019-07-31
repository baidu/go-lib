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

package web_params

import (
	"testing"
)

func TestParamsValueGet(t *testing.T) {
	params := make(map[string][]string)

	params["format"] = []string{"json", "kv"}

	value, err := ParamsValueGet(params, "form")
	if err == nil {
		t.Error("err in ParamsValueGet(), should no value for 'form'")
	}

	value, err = ParamsValueGet(params, "format")
	if value != "json" {
		t.Error("err in ParamsValueGet(), value should be 'json'")
	}
}

func TestParamsMultiValueGet(t *testing.T) {
	params := make(map[string][]string)

	params["format"] = []string{"json", "kv"}

	value, err := ParamsMultiValueGet(params, "form")
	if err == nil {
		t.Error("err in ParamsMultiValueGet(), should no value for 'form'")
	}

	value, err = ParamsMultiValueGet(params, "format")
	if value == nil {
		t.Error("err in ParamsMultiValueGet(), value should not be nil")
		return
	}
	if len(value) != 2 || value[0] != "json" || value[1] != "kv" {
		t.Error("err in ParamsMultiValueGet(), err in value for 'format'")
	}
}
