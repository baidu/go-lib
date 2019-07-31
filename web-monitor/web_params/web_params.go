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
	"errors"
)

// ParamsValueGet gets one (the first) value for given key in params
func ParamsValueGet(params map[string][]string, key string) (string, error) {
	values := params[key]

	if values == nil || len(values) == 0 {
		return "", errors.New("key not exist")
	}

	return values[0], nil
}

// ParamsMultiValueGet gets values for given key in params
func ParamsMultiValueGet(params map[string][]string, key string) ([]string, error) {
	values := params[key]

	if values == nil || len(values) == 0 {
		return nil, errors.New("key not exist")
	}

	return values, nil
}
