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

// Package log4go provides level-based and highly configurable logging.
//
// Enhanced Logging
//
// This is inspired by the logging functionality in Java.  Essentially, you create a Logger
// object and create output filters for it.  You can send whatever you want to the Logger,
// and it will filter that based on your settings and send it to the outputs.  This way, you
// can put as much debug code in your program as you want, and when you're done you can filter
// out the mundane messages so only the important ones show up.
//
// Utility functions are provided to make life easier. Here is some example code to get started:
//
// log := log4go.NewLogger()
// log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())
// log.AddFilter("log",    log4go.FINE,  log4go.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// The first two lines can be combined with the utility NewDefaultLogger:
//
// log := log4go.NewDefaultLogger(log4go.DEBUG)
// log.AddFilter("log",    log4go.FINE,  log4go.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// Usage notes:
// - The ConsoleLogWriter does not display the source of the message to standard
//   output, but the FileLogWriter does.
// - The utility functions (Info, Debug, Warn, etc) derive their source from the
//   calling function, and this incurs extra overhead.
//
// Changes from 2.0:
// - The external interface has remained mostly stable, but a lot of the
//   internals have been changed, so if you depended on any of this or created
//   your own LogWriter, then you will probably have to update your code.  In
//   particular, Logger is now a map and ConsoleLogWriter is now a channel
//   behind-the-scenes, and the LogWrite method no longer has return values.
//
// Future work: (please let me know if you think I should work on any of these particularly)
// - Log file rotation
// - Logging configuration files ala log4j
// - Have the ability to remove filters?
// - Have GetInfoChannel, GetDebugChannel, etc return a chan string that allows
//   for another method of logging
// - Add an XML filter type
package log4go

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (log Logger) LogWithContext(ctx context.Context, lvl LevelType, source, message string) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Make the log record
	rec := &LogRecord{
		Ctx:     ctx,
		Level:   lvl,
		Created: time.Now(),
		Source:  source,
		Message: message,
		Binary:  nil,
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Logf logs a formatted log message at the given log level, using the caller as
// its source.
func (log Logger) LogfWithContext(ctx context.Context, lvl LevelType, format string, args ...interface{}) {
	log.intLogf(ctx, lvl, format, args...)
}

// Logc logs a string returned by the closure at the given log level, using the caller as
// its source.  If no log message would be written, the closure is never called.
func (log Logger) LogcWithContext(ctx context.Context, lvl LevelType, closure func() string) {
	log.intLogc(ctx, lvl, closure)
}

// FinestWithContext logs a message at the finest log level.
// See Debug for an explanation of the arguments.
func (log Logger) FinestWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// FineWithContext logs a message at the fine log level.
// See Debug for an explanation of the arguments.
func (log Logger) FineWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// DebugWithContext is a utility method for debug log messages.
// The behavior of Debug depends on the first argument:
// - arg0 is a string
//   When given a string as the first argument, this behaves like Logf but with
//   the DEBUG log level: the first argument is interpreted as a format for the
//   latter arguments.
// - arg0 is a func()string
//   When given a closure of type func()string, this logs the string returned by
//   the closure iff it will be logged.  The closure runs at most one time.
// - arg0 is interface{}
//   When given anything else, the log message will be each of the arguments
//   formatted with %v and separated by spaces (ala Sprint).
func (log Logger) DebugWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// TraceWithContext logs a message at the trace log level.
// See Debug for an explanation of the arguments.
func (log Logger) TraceWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// InfoWithContext logs a message at the info log level.
// See Debug for an explanation of the arguments.
func (log Logger) InfoWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(ctx, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(ctx, lvl, first)
	case []byte:
		// Log the binary log message
		log.intLogb(ctx, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(ctx, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// WarnWithContext logs a message at the warning log level and returns the formatted error.
// At the warning level and higher, there is no performance benefit if the
// message is not actually logged, because all formats are processed and all
// closures are executed to format the error message.
// See Debug for further explanation of the arguments.
func (log Logger) WarnWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(ctx, lvl, msg)
	return errors.New(msg)
}

// ErrorWithContext logs a message at the error log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (log Logger) ErrorWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(ctx, lvl, msg)
	return errors.New(msg)
}

// CriticalWithContext logs a message at the critical log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (log Logger) CriticalWithContext(ctx context.Context, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(ctx, lvl, msg)
	return errors.New(msg)
}
