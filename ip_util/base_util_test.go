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
    "net"
    "testing"
)

func TestCopyBoolSlice(t *testing.T) {
    boolSlice := []bool{false, false, false, true}
    copiedSlice := copyBoolSlice(boolSlice)
    if &copiedSlice[0] == &boolSlice[0] {
        t.Errorf("the address of original bool slice and copied should be different but not")
        return
    }
}

func TestGenMinAddress(t *testing.T) {
    ip := []bool{false, false, false, true, false, true, true, true}
    minAddress := genMinAddress(ip, 3)
    for i := 0; i < len(ip); i++ {
        if i <= 3 {
            if minAddress[i] != ip[i] {
                t.Errorf("minAddress[%d] should be equal to ip[%d] but not", i, i)
                return
            }
        } else {
            if minAddress[i] != false {
                t.Errorf("minAddress[%d] shoule be false but not", i)
                return
            }
        }
    }
}

func TestGenMaxAddress(t *testing.T) {
    ip := []bool{false, false, false, true, false, true, true, true}
    maxAddress := genMaxAddress(ip, 3)
    for i := 0; i < len(ip); i++ {
        if i <= 3 {
            if maxAddress[i] != ip[i] {
                t.Errorf("maxAddress[%d] should be equal to ip[%d] but not", i, i)
                return
            }
        } else {
            if maxAddress[i] != true {
                t.Errorf("maxAddress[%d] should be true but not", i)
                return
            }
        }
    }
}

func TestLowerEqual(t *testing.T) {
    var testCases = []struct {
        ip1  []bool
        ip2  []bool
        want bool
    }{
        {
            []bool{false, false, false, true},
            []bool{false, false, true, false},
            true,
        },

        {
            []bool{false, false, true, true},
            []bool{false, false, false, true},
            false,
        },

        {
            []bool{false, false, true, true},
            []bool{false, false, true, true},
            true,
        },
    }

    for index, testCase := range testCases {
        if lowerEqual(testCase.ip1, testCase.ip2) != testCase.want {
            t.Errorf("run testCases[%d] failed", index)
            return
        }
    }
}

func TestConvertByteToBoolSlice(t *testing.T) {
    var testCases = []struct {
        input byte
        want  []bool
    }{
        {
            0,
            []bool{false, false, false, false, false, false, false, false},
        },
        {
            2,
            []bool{false, false, false, false, false, false, true, false},
        },
        {
            7,
            []bool{false, false, false, false, false, true, true, true},
        },
    }

    for index, testCase := range testCases {
        result := convertByteToBoolSlice(testCase.input)
        for i := 0; i < len(testCase.want); i++ {
            if result[i] != testCase.want[i] {
                t.Errorf("run testCases[%d] failed", index)
                return
            }
        }
    }
}

func TestConvertIpToBoolSlice(t *testing.T) {
    ip := net.ParseIP("192.168.1.0")
    if ip == nil {
        t.Errorf("192.168.1.0: invalid ip format")
        return
    }
    ip = ip[12:]
    ipBinary := convertIpToBoolSlice(ip)
    byteArr := []byte{192, 168, 1, 0}
    base := 0
    for i := 0; i < len(byteArr); i++ {
        ipBinarySection := convertByteToBoolSlice(byteArr[i])
        for j := base; j < base+8; j++ {
            if ipBinary[j] != ipBinarySection[j-base] {
                t.Errorf("ipBinary[%d] should be equal to ipBinarySection[%d] but not", j, j-base)
                return
            }
        }
        base += 8
    }
}

func TestFirstDiffPos(t *testing.T) {
    boolSlice1 := []bool{false, false, false, false}
    boolSlice2 := []bool{false, false, true, false}
    pos, err := firstDiffPos(boolSlice1, boolSlice2)
    if err != nil {
        t.Errorf("firstDiffPos(): %s", err.Error())
        return
    }
    if pos != 2 {
        t.Errorf("firstDiffPos() should return 2 but not")
        return
    }

    boolSlice1 = []bool{false, false, true, true}
    boolSlice2 = []bool{false, false, true, true}
    pos, err = firstDiffPos(boolSlice1, boolSlice2)
    if err == nil {
        t.Errorf("boolSlice1 is equal to boolSlice2 and firstDiffPos() should return error but not")
        return
    }
}

func TestConvertIpSectionToInt(t *testing.T) {
    ipSection := []bool{false, false, false, false, true, true, true, true}
    intValue := convertIPSectionToInt(ipSection)
    if intValue != 15 {
        t.Errorf("convertIPSectionToInt() should return 15 but not")
        return
    }
}

func TestConvertIPv6SectionToString(t *testing.T) {
    testCases := []struct {
        input []bool
        want  string
    }{
        {
            []bool{false, false, false, false, false, false, false, false,
                false, false, false, false, false, false, false, false},
            "0",
        },
        {
            []bool{false, false, false, false, false, false, true, true,
                false, true, true, true, true, true, true, true},
            "37f",
        },
    }

    for index, testCase := range testCases {
        get := convertIPv6SectionToString(testCase.input)
        if get != testCase.want {
            t.Errorf("run testCases[%d] failed, want %s but get %s", index, testCase.want, get)
            return
        }
    }
}

func TestConvertIPv4SectionToString(t *testing.T) {
    ipSection := []bool{true, true, false, false, false, false, false, false}
    if convertIPv4SectionToString(ipSection) != "192" {
        t.Errorf("convertIPv4SectionToString should return 192 but not")
        return
    }
}

func TestGenIPv6Cidr(t *testing.T) {
    ip := net.ParseIP("0:0:0:0:0:0:ff:0")
    if ip == nil {
        t.Errorf("0:0:0:0:0:0:ff:0: invalid ip format")
        return
    }
    ipBinary := convertIpToBoolSlice(ip)
    cidr, err := genCidr(ipBinary, 112)
    if err != nil {
        t.Errorf("genCidr(): %s", err.Error())
        return
    }
    if cidr != "::ff:0/112" {
        t.Errorf("genCidr() should return ::ff:0/112 but return %s", cidr)
        return
    }
}

func TestGenIPv4Cidr(t *testing.T) {
    ip := net.ParseIP("192.168.1.0")
    if ip == nil {
        t.Errorf("192.168.1.0: invalid ip format")
        return
    }
    ip = ip[12:]
    ipBinary := convertIpToBoolSlice(ip)
    cidr, err := genIPv4Cidr(ipBinary, 24)
    if err != nil {
        t.Errorf("genCidr(): %s", err.Error())
        return
    }
    if cidr != "192.168.1.0/24" {
        t.Errorf("genCidr() should return 192.168.1.0/24 but return %s", cidr)
        return
    }
}

func TestGenCidr(t *testing.T) {
    ip := net.ParseIP("192.168.1.0")
    if ip == nil {
        t.Errorf("192.168.1.0: invalid ip format")
        return
    }
    ip = ip[12:]
    ipBinary := convertIpToBoolSlice(ip)
    cidr, err := genCidr(ipBinary, 24)
    if err != nil {
        t.Errorf("genCidr(): %s", err.Error())
        return
    }
    if cidr != "192.168.1.0/24" {
        t.Errorf("genCidr() should return 192.168.1.0/24 but return %s", cidr)
        return
    }

    ip = net.ParseIP("0:0:0:0:0:0:ff:0")
    if ip == nil {
        t.Errorf("0:0:0:0:0:0:ff:0: invalid ip format")
        return
    }
    ipBinary = convertIpToBoolSlice(ip)
    cidr, err = genCidr(ipBinary, 112)
    if err != nil {
        t.Errorf("genCidr(): %s", err.Error())
        return
    }
    if cidr != "::ff:0/112" {
        t.Errorf("genCidr() should return ::ff:0/112 but return %s", cidr)
        return
    }
}