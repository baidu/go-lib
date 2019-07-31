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

package kv_encode

import (
	"fmt"
	"strings"
	"testing"
)

type testData struct {
	a int
	b string
	c int32
}

func TestEncode(t *testing.T) {
	var data testData

	data.a = 123
	data.b = "456"
	data.c = 789

	buf, err := Encode(data)

	if err != nil {
		errStr := fmt.Sprintf("err in Encode():%s", err.Error())
		t.Error(errStr)
		return
	}

	str := string(buf)
	str = strings.TrimSuffix(str, "\n")
	strs := strings.Split(str, "\n")

	strMap := map[string]bool{
		"a:123": true,
		"b:456": true,
		"c:789": true,
	}

	for _, str = range strs {
		_, ok := strMap[str]
		if !ok {
			t.Error("err in Encode(): result is not expected")
			return
		}

		delete(strMap, str)
	}
}
