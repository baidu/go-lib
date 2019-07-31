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

// Package log is an encapsulation for log4go

/*
Usage:
    import "github.com/baidu/go-lib/log"

    // Two log files will be generated in ./log:
    // test.log, and test.wf.log(for log > warn)
    // The log will rotate, and there is support for backup count
    log.Init("test", "INFO", "./log", true, "midnight", 5)

    log.Logger.Warn("warn msg")
    log.Logger.Info("info msg")

    // it is required, to work around bug of log4go
    time.Sleep(100 * time.Millisecond)
*/
package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

import "github.com/baidu/go-lib/log/log4go"

// Logger is global logger
var Logger log4go.Logger
var initialized bool = false
var mutex sync.Mutex

// logDirCreate checks and creates dir if nonexist
func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(logDir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// filenameGen generates filename
func filenameGen(progName, logDir string, isErrLog bool) string {
	/* remove the last '/'  */
	strings.TrimSuffix(logDir, "/")

	var fileName string
	if isErrLog {
		/* for log file of warning, error, critical  */
		fileName = filepath.Join(logDir, progName+".wf.log")
	} else {
		/* for log file of all log  */
		fileName = filepath.Join(logDir, progName+".log")
	}

	return fileName
}

// stringToLevel converts level in string to log4go level
func stringToLevel(str string) log4go.LevelType {
	var level log4go.LevelType

	str = strings.ToUpper(str)

	switch str {
	case "DEBUG":
		level = log4go.DEBUG
	case "TRACE":
		level = log4go.TRACE
	case "INFO":
		level = log4go.INFO
	case "WARNING":
		level = log4go.WARNING
	case "ERROR":
		level = log4go.ERROR
	case "CRITICAL":
		level = log4go.CRITICAL
	default:
		level = log4go.INFO
	}
	return level
}

// Init initializes log lib
//
// PARAMS:
//   - progName: program name. Name of log file will be progName.log
//   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
//   - logDir: directory for log. It will be created if noexist
//   - hasStdOut: whether to have stdout output
//   - when:
//       "M", minute
//       "H", hour
//       "D", day
//       "MIDNIGHT", roll over at midnight
//   - backupCount: If backupCount is > 0, when rollover is done, no more than
//       backupCount files are kept - the oldest ones are deleted.
func Init(progName string, levelStr string, logDir string,
	hasStdOut bool, when string, backupCount int) error {
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return errors.New("Initialized Already")
	}

	var err error
	Logger, err = Create(progName, levelStr, logDir, hasStdOut, when, backupCount)
	if err != nil {
		return err
	}

	initialized = true
	return nil
}

// Create creates log lib
//
// PARAMS:
//   - progName: program name. Name of log file will be progName.log
//   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
//   - logDir: directory for log. It will be created if noexist
//   - hasStdOut: whether to have stdout output
//   - when:
//       "M", minute
//       "H", hour
//       "D", day
//       "MIDNIGHT", roll over at midnight
//   - backupCount: If backupCount is > 0, when rollover is done, no more than
//       backupCount files are kept - the oldest ones are deleted.
func Create(progName string, levelStr string, logDir string,
	hasStdOut bool, when string, backupCount int) (log4go.Logger, error) {
	/* check when   */
	if !log4go.WhenIsValid(when) {
		return nil, fmt.Errorf("invalid value of when: %s", when)
	}

	/* check, and create dir if nonexist    */
	if err := logDirCreate(logDir); err != nil {
		log4go.Error("Init(), in logDirCreate(%s)", logDir)
		return nil, err
	}

	/* convert level from string to log4go level    */
	level := stringToLevel(levelStr)

	/* create logger    */
	logger := make(log4go.Logger)

	/* create writer for stdout */
	if hasStdOut {
		logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	/* create file writer for all log   */
	fileName := filenameGen(progName, logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return nil, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	logWriter.SetFormat(log4go.LogFormat)
	logger.AddFilter("log", level, logWriter)

	/* create file writer for warning and fatal log */
	fileNameWf := filenameGen(progName, logDir, true)
	logWriter = log4go.NewTimeFileLogWriter(fileNameWf, when, backupCount)
	if logWriter == nil {
		return nil, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileNameWf)
	}
	logWriter.SetFormat(log4go.LogFormat)
	logger.AddFilter("log_wf", log4go.WARNING, logWriter)

	return logger, nil
}

// InitWithLogSvr initializes log lib with remote log server
//
// PARAMS:
//   - progName: program name.
//   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
//   - loggerName: logger name
//   - network: using "udp" or "unixgram"
//   - svrAddr: remote unix sock address for all logger
//   - svrAddrWf: remote unix sock address for warn/fatal logger
//                If svrAddrWf is empty string, no warn/fatal logger will be created.
//   - hasStdOut: whether to have stdout output
func InitWithLogSvr(progName string, levelStr string, loggerName string,
	network string, svrAddr string, svrAddrWf string,
	hasStdOut bool) error {
	if initialized {
		return errors.New("Initialized Already")
	}

	/* convert level from string to log4go level    */
	level := stringToLevel(levelStr)

	/* create logger    */
	Logger = make(log4go.Logger)

	/* create writer for stdout */
	if hasStdOut {
		Logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	/* create file writer for all log   */
	name := fmt.Sprintf("%s_%s", progName, loggerName)

	logWriter := log4go.NewPacketWriter(name, network, svrAddr, log4go.LogFormat)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewPacketWriter(%s)", name)
	}
	Logger.AddFilter("log", level, logWriter)

	if len(svrAddrWf) > 0 {
		/* create file writer for warning and fatal log */
		logWriterWf := log4go.NewPacketWriter(name+".wf", network, svrAddrWf, log4go.LogFormat)
		if logWriterWf == nil {
			return fmt.Errorf("error in log4go.NewPacketWriter(%s, %s)",
				name, svrAddr)
		}
		Logger.AddFilter("log_wf", log4go.WARNING, logWriterWf)
	}

	initialized = true
	return nil
}
