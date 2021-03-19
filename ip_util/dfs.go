// Copyright (c) 2021 Baidu, Inc.
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

package ip_util

var cidrList []string

// depth-first search complete cidr list to represent ip range [ipBinaryBegin, ipBinaryEnd]
func dfs(ip []bool, pos int, isFirst bool, ipBinaryBegin, ipBinaryEnd []bool) error {
    if pos >= len(ipBinaryBegin) {
        return nil
    }

    var minAddress []bool
    var maxAddress []bool
    prefixLength := pos + 1
    if isFirst {
        minAddress = genMinAddress(ip, pos)
        maxAddress = genMaxAddress(ip, pos)
        if lowerEqual(ipBinaryBegin, minAddress) && lowerEqual(maxAddress, ipBinaryEnd) {
            cidr, err := genCidr(ip, prefixLength)
            if err != nil {
                return err
            }
            cidrList = append(cidrList, cidr)
            return nil
        }
        return dfs(ip, pos + 1, false, ipBinaryBegin, ipBinaryEnd)
    }

    ip[pos] = false
    minAddress = genMinAddress(ip, pos)
    maxAddress = genMaxAddress(ip, pos)
    if lowerEqual(ipBinaryBegin, minAddress) && lowerEqual(maxAddress, ipBinaryEnd) {
        cidr, err := genCidr(ip, prefixLength)
        if err != nil {
            return err
        }
        cidrList = append(cidrList, cidr)
    } else {
        err := dfs(ip, pos + 1, false, ipBinaryBegin, ipBinaryEnd)
        if err != nil {
            return nil
        }
    }

    ip[pos] = true
    minAddress = genMinAddress(ip, pos)
    maxAddress = genMaxAddress(ip, pos)
    if lowerEqual(ipBinaryBegin, minAddress) && lowerEqual(maxAddress, ipBinaryEnd) {
        cidr, err := genCidr(ip, prefixLength)
        if err != nil {
            return err
        }
        cidrList = append(cidrList, cidr)
    } else {
        err := dfs(ip, pos + 1, false, ipBinaryBegin, ipBinaryEnd)
        if err != nil {
            return err
        }
    }

    return nil
}
