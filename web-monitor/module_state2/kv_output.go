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

package module_state2

import (
	"fmt"
	"strings"
	"unicode"
)

// keyGen generate key for key-value output
// Params:
//  - key: the original key
//  - keyPrefix: e.g., "mod_header"
//  - programName: e.g., "bfe"
//  - withProgramName: whether program name should be included in the result
//
// Returns:
//  final key, e.g., "mod_header_ERR_PB_SEEK", "bfe.mod_header_ERR_PB_SEEK"(with programName)
func keyGen(key string, keyPrefix string, programName string, withProgramName bool) string {
	if programName != "" && withProgramName {
		if keyPrefix == "" {
			return fmt.Sprintf("%s.%s", programName, key)
		}

		return fmt.Sprintf("%s.%s_%s", programName, keyPrefix, key)
	} else {
		if keyPrefix == "" {
			return key
		}

		return fmt.Sprintf("%s_%s", keyPrefix, key)
	}
}

// escapeKey replaces unsupport character with "_"
// Some monitor system only support letter, num, ".", "_" for key
// 
// Params:
//     - originKey: original key
// 
// Returns:
//     escaped string
func escapeKey(originKey string) string {
	finalKey := originKey
	for _, str := range originKey {
		switch {
		case unicode.IsLetter(str):
		case unicode.IsNumber(str):
		case byte(str) == '-':
		case byte(str) == '_':
		case byte(str) == '.':
		default:
			finalKey = strings.Replace(finalKey, string(str), "_", -1)
		}
	}
	return finalKey
}

// KeyGen generate and escape key for key-value output
//
// Params:
//  - key: the original key
//  - keyPrefix: e.g., "mod_header"
//  - programName: e.g., "bfe"
//  - withProgramName: whether program name should be included in the result
//
// Returns:
//  final key, e.g., "mod_header_ERR_PB_SEEK", "bfe.mod_header_ERR_PB_SEEK"(with programName)
func KeyGen(key string, keyPrefix string, programName string, withProgramName bool) string {
	finalKey := keyGen(key, keyPrefix, programName, withProgramName)
	return escapeKey(finalKey)
}
