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

// pool implements pool for binary data

package log4go

import (
	"errors"
	"sync"
)

// buffer pool for binary data
var buf4kPool sync.Pool
var buf16kPool sync.Pool
var buf64kPool sync.Pool

var (
	ErrTooLarge = errors.New("required slice size exceed 64k")
)

// NewBuffer gets proper []byte from pool
// if size > 16K, return ErrTooLarge
func NewBuffer(size int) ([]byte, error) {
	var pool *sync.Pool

	// return buffer size
	originSize := size

	if size <= 4096 {
		size = 4096
		pool = &buf4kPool
	} else if size <= 16*1024 {
		size = 16 * 1024
		pool = &buf16kPool
	} else if size <= 64*1024 {
		size = 64 * 1024
		pool = &buf64kPool
	} else {
		// if message is larger than 64K, return err
		return nil, ErrTooLarge
	}

	if v := pool.Get(); v != nil {
		return v.([]byte)[:originSize], nil
	}

	return make([]byte, size)[:originSize], nil
}

func putBuffer(b []byte) {
	b = b[:cap(b)]
	if cap(b) == 4096 {
		buf4kPool.Put(b)
	}
	if cap(b) == 16*1024 {
		buf16kPool.Put(b)
	}
	if cap(b) == 64*1024 {
		buf64kPool.Put(b)
	}
}
