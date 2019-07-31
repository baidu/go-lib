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
    "testing"
)

// test WhenIsValid()
func TestWhenIsValid(t *testing.T) {
    if ! WhenIsValid("midNIGHT") {
        t.Error("err in WhenIsValid('midNIGHT')")
    }

    if WhenIsValid("mid-night") {
        t.Error("err in WhenIsValid('mid-night')")
    }

    if !WhenIsValid("m") {
        t.Error("err in WhenIsValid('m')")
    }

    if !WhenIsValid("H") {
        t.Error("err in WhenIsValid('H')")
    }

    if !WhenIsValid("d") {
        t.Error("err in WhenIsValid('H')")
    }
}