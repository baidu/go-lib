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

// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

func CrashWithContext(ctx context.Context, args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ctx, CRITICAL, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Logs the given message and crashes the program
func CrashfWithContext(ctx context.Context, format string, args ...interface{}) {
	Global.intLogf(ctx, CRITICAL, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Compatibility with `log`
func ExitWithContext(ctx context.Context, args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ctx, ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func ExitfWithContext(ctx context.Context, format string, args ...interface{}) {
	Global.intLogf(ctx, ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func StderrWithContext(ctx context.Context, args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ctx, ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func StderrfWithContext(ctx context.Context, format string, args ...interface{}) {
	Global.intLogf(ctx, ERROR, format, args...)
}

// Compatibility with `log`
func StdoutWithContext(ctx context.Context, args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ctx, INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func StdoutfWithContext(ctx context.Context, format string, args ...interface{}) {
	Global.intLogf(ctx, INFO, format, args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func LogWithContext(ctx context.Context, lvl LevelType, source, message string) {
	Global.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func LogfWithContext(ctx context.Context, lvl LevelType, format string, args ...interface{}) {
	Global.intLogf(ctx, lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func LogcWithContext(ctx context.Context, lvl LevelType, closure func() string) {
	Global.intLogc(ctx, lvl, closure)
}

// Utility for finest log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Finest
func FinestWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for fine log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Fine
func FineWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func DebugWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func TraceWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func InfoWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func WarnWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(ctx, lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func ErrorWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(ctx, lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical
func CriticalWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(ctx, lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(ctx, lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(ctx, lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}
