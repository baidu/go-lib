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

import (
    "fmt"
    "net"
    "strconv"
)

// copy bool slice
func copyBoolSlice(ip []bool) []bool {
    cp := make([]bool, len(ip))
    copy(cp, ip)
    return cp
}

// generate min address under the condition that we can not modify the bit of ip before pos
func genMinAddress(ip []bool, pos int) []bool {
    ip = copyBoolSlice(ip)
    for i := pos + 1; i < len(ip); i++ {
        ip[i] = false
    }
    return ip
}

// generate max address under the condition that we can not modify the bit of ip before pos
func genMaxAddress(ip []bool, pos int) []bool {
    ip = copyBoolSlice(ip)
    for i := pos + 1; i < len(ip); i++ {
        ip[i] = true
    }
    return ip
}

// just like compare numbers, return whether section1 <= section2
func lowerEqual(section1, section2 []bool) bool {
    n := len(section1)
    for i := 0; i < n; i++ {
        if !section1[i] && section2[i] {
            return true
        }
        if section1[i] && !section2[i] {
            return false
        }
    }
    return true
}

// convert byte to bool slice
func convertByteToBoolSlice(section byte) []bool {
    var boolSlice []bool
    var p byte = 128
    for p > 0 {
        if section & p == 0 {
            boolSlice = append(boolSlice, false)
        } else {
            boolSlice = append(boolSlice, true)
        }
        p >>= 1
    }
    return boolSlice
}

// convert ip to bool slice
func convertIpToBoolSlice(ip []byte) []bool {
    var bitSlice []bool
    for _, section := range ip {
        boolSlice := convertByteToBoolSlice(section)
        for _, bit := range boolSlice {
            bitSlice = append(bitSlice, bit)
        }
    }
    return bitSlice
}

// get the position of first different bit
func firstDiffPos(ipBinaryBegin, ipBinaryEnd []bool) (int, error) {
    for i := 0; i < len(ipBinaryBegin); i++ {
        if ipBinaryBegin[i] != ipBinaryEnd[i] {
            return i, nil
        }
    }
    return 0, fmt.Errorf("no diff")
}

// convert ip section to int
func convertIPSectionToInt(ipSection []bool) int {
    decimal := 0
    base := 1
    for i := len(ipSection) - 1; i >= 0; i-- {
        if ipSection[i] {
            decimal += base
        }
        base <<= 1
    }
    return decimal
}

// convert IPv6 section to string
func convertIPv6SectionToString(ipSection []bool) string {
    var strBytes []byte
    // 二进制转int
    decimal := convertIPSectionToInt(ipSection)
    // int转16进制
    if decimal == 0{
        return "0"
    }
    for decimal > 0 {
        mod := decimal % 16
        if mod >= 10 {
            strBytes = append(strBytes, 'a' + byte(mod - 10))
        } else {
            strBytes = append(strBytes, '0' + byte(mod))
        }
        decimal >>= 4
    }
    // reverse
    n := len(strBytes)
    for i := 0; i < n / 2; i++ {
        strBytes[i], strBytes[n - i - 1] = strBytes[n - i - 1], strBytes[i]
    }
    return string(strBytes)
}

// convert IPv4 section to string
func convertIPv4SectionToString(ipSection []bool) string {
    return strconv.Itoa(convertIPSectionToInt(ipSection))
}

// generate IPv6 cidr
func genIPv6Cidr(ip []bool, prefixLength int) (string, error) {
    cidr := ""
    for i := 0; i < 8; i++ {
        if i > 0 {
            cidr += ":"
        }
        cidr += convertIPv6SectionToString(ip[i * 16:(i + 1) * 16])
    }
    cidr += "/" + strconv.Itoa(prefixLength)
    _, ipNet, err := net.ParseCIDR(cidr)
    if err != nil {
        return "", err
    }
    return ipNet.String(), nil
}

// generate IPv4 cidr
func genIPv4Cidr(ip []bool, prefixLength int) (string, error) {
    cidr := ""
    for i := 0; i < 4; i++ {
        if i > 0 {
            cidr += "."
        }
        cidr += convertIPv4SectionToString(ip[i * 8:(i + 1) * 8])
    }
    cidr += "/" + strconv.Itoa(prefixLength)
    _, ipNet, err := net.ParseCIDR(cidr)
    if err != nil {
        return "", err
    }
    return ipNet.String(), nil
}

// generate cidr
func genCidr(ip []bool, prefixLength int) (string, error) {
    if len(ip) == 128 {
        return genIPv6Cidr(ip, prefixLength)
    }
    return genIPv4Cidr(ip, prefixLength)
}