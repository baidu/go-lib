// Copyright (c) 2020 Baidu, Inc.
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
/*
DESCRIPTION
*/
package queue

import (
    "fmt"
    "testing"
    "time"
)

func consumer(queue *Queue) {
    for {
        number := queue.Remove()
        fmt.Println("Read from queue: ", number)
    }
}

func producer(queue *Queue, t *testing.T) {
    for i := 0; i < 10; i = i + 1 {
        retVal := queue.Append(i)
        if retVal != nil {
            t.Error("queue.Append() should return nil")
        }        
        fmt.Println("write to queue: ", i)
    }
}

func TestSendQueue(t *testing.T) {
    var queue Queue
    queue.Init()
        
    go consumer(&queue)
    go producer(&queue, t)

    time.Sleep(2 * time.Second)
}

func TestQueueIsFull(t *testing.T) {
    var queue Queue
    queue.Init()
    queue.SetMaxLen(3)
        
    for i := 0; i < 10; i = i + 1 {
        retVal := queue.Append(i)
        
        if i < 3 {
            if retVal != nil {
                t.Error("queue.Append() should return nil")
            }
        } else {
            if retVal == nil {
                t.Error("queue.Append() should return error")
            }            
        }
    }    
}

