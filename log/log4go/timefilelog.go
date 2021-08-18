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

// timefilelog implements enhanced timed file log for go
/*
This file is modified from filelog.go

It supports:
- Split log file by day, hour, minite
- Suffix of log file reflect time of logging
- Support backupCount
*/
package log4go

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

import (
	strftime "github.com/jehiah/go-strftime"
)

const (
	MIDNIGHT = 24 * 60 * 60 /* number of seconds in a day */
	NEXTHOUR = 60 * 60      /* number of seconds in a hour */
)

// This log writer sends output to a file
type TimeFileLogWriter struct {
	LogCloser //for Elegant exit

	rec chan *LogRecord

	// The opened file
	filename     string
	baseFilename string // abs path
	file         *os.File

	// The logging format
	format string

	when        string // 'D', 'H', 'M', "MIDNIGHT", "NEXTHOUR"
	backupCount int    // If backupCount is > 0, when rollover is done,
	// no more than backupCount files are kept

	interval   int64
	suffix     string         // suffix of log file
	fileFilter *regexp.Regexp // for removing old log files

	rolloverAt int64 // time.Unix()
}

// WhenIsValid checks whether value of when is valid
func WhenIsValid(when string) bool {
	switch strings.ToUpper(when) {
	case "MIDNIGHT", "NEXTHOUR", "M", "H", "D":
		return true
	default:
		return false
	}
}

// This is the FileLogWriter's output method
func (w *TimeFileLogWriter) LogWrite(rec *LogRecord) {
	if !LogWithBlocking {
		if len(w.rec) >= LogBufferLength {
			return
		}
	}

	w.rec <- rec
}

// Close waits for dump all log and close chan
func (w *TimeFileLogWriter) Close() {
	w.WaitForEnd(w.rec)
	close(w.rec)
}

func (w *TimeFileLogWriter) computeRollover(currTime time.Time) int64 {
	var result int64

	if w.when == "MIDNIGHT" {
		t := currTime.Local()
		/* r is the number of seconds left between now and midnight */
		r := MIDNIGHT - ((t.Hour()*60+t.Minute())*60 + t.Second())
		result = currTime.Unix() + int64(r)
	} else if w.when == "NEXTHOUR" {
		t := currTime.Local()
		/* r is the number of seconds left between now and the next hour */
		r := NEXTHOUR - (t.Minute()*60 + t.Second())
		result = currTime.Unix() + int64(r)
	} else {
		result = currTime.Unix() + w.interval
	}
	return result
}

// prepare prepares according to "when"
func (w *TimeFileLogWriter) prepare() {
	var regRule string

	switch w.when {
	case "M":
		w.interval = 60
		w.suffix = "%Y-%m-%d_%H-%M"
		regRule = `^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}$`
	case "H", "NEXTHOUR":
		w.interval = 60 * 60
		w.suffix = "%Y-%m-%d_%H"
		regRule = `^\d{4}-\d{2}-\d{2}_\d{2}$`
	case "D", "MIDNIGHT":
		w.interval = 60 * 60 * 24
		w.suffix = "%Y-%m-%d"
		regRule = `^\d{4}-\d{2}-\d{2}$`
	default:
		// default is "D"
		w.interval = 60 * 60 * 24
		w.suffix = "%Y-%m-%d"
		regRule = `^\d{4}-\d{2}-\d{2}$`
	}
	w.fileFilter = regexp.MustCompile(regRule)

	fInfo, err := os.Stat(w.filename)

	var t time.Time
	if err == nil {
		t = fInfo.ModTime()
	} else {
		t = time.Now()
	}

	w.rolloverAt = w.computeRollover(t)
}

func (w *TimeFileLogWriter) shouldRollover() bool {
	t := time.Now().Unix()

	if t >= w.rolloverAt {
		return true
	} else {
		return false
	}
}

// NewTimeFileLogWriter creates a new TimeFileLogWriter
//
// PARAMS:
//   - fname: name of log file
//   - when:
//       "M", minute
//       "H", hour
//       "D", day
//       "MIDNIGHT", roll over at midnight
//       "NEXTHOUR", roll over at sharp next hour
//   - backupCount: If backupCount is > 0, when rollover is done, no more than
//       backupCount files are kept - the oldest ones are deleted.
func NewTimeFileLogWriter(fname string, when string, backupCount int) *TimeFileLogWriter {
	// check value of when is valid
	if !WhenIsValid(when) {
		fmt.Fprintf(os.Stderr, "NewTimeFileLogWriter(%q): invalid value of when:%s \n",
			fname, when)
		return nil
	}

	// change when to upper
	when = strings.ToUpper(when)

	// create TimeFileLogWriter
	w := &TimeFileLogWriter{
		rec:         make(chan *LogRecord, LogBufferLength),
		filename:    fname,
		format:      "[%D %T] [%L] (%S) %M",
		when:        when,
		backupCount: backupCount,
	}

	// add w to collection of all writers
	writersInfo = append(writersInfo, w)

	//init LogCloser
	w.LogCloserInit()

	// get abs path
	if path, err := filepath.Abs(fname); err != nil {
		fmt.Fprintf(os.Stderr, "NewTimeFileLogWriter(%q): %s\n", w.filename, err)
		return nil
	} else {
		w.baseFilename = path
	}

	// prepare for w.interval, w.suffix and w.fileFilter
	w.prepare()

	// open the file for the first time
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "NewTimeFileLogWriter(%q): %s\n", w.filename, err)
		return nil
	}

	go func() {
		defer func() {
			if w.file != nil {
				w.file.Close()
			}
		}()

		for rec := range w.rec {
			if w.EndNotify(rec) {
				return
			}

			if w.shouldRollover() {
				if err := w.intRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "NewTimeFileLogWriter(%q): %s\n", w.filename, err)
					return
				}
			}

			// Perform the write
			var err error
			if rec.Binary != nil {
				_, err = w.file.Write(rec.Binary)
			} else {
				_, err = fmt.Fprint(w.file, FormatLogRecord(w.format, rec))
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "NewTimeFileLogWriter(%q): %s\n", w.filename, err)
				return
			}
		}
	}()

	return w
}

// getFilesToDelete determines the files to delete when rolling over
func (w *TimeFileLogWriter) getFilesToDelete() []string {
	dirName := filepath.Dir(w.baseFilename)
	baseName := filepath.Base(w.baseFilename)

	result := []string{}

	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
		return result
	}

	prefix := baseName + "."
	plen := len(prefix)

	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()
		if len(fileName) >= plen {
			if fileName[:plen] == prefix {
				suffix := fileName[plen:]
				if w.fileFilter.MatchString(suffix) {
					result = append(result, filepath.Join(dirName, fileName))
				}
			}
		}
	}

	sort.Sort(sort.StringSlice(result))

	if len(result) < w.backupCount {
		result = result[0:0]
	} else {
		result = result[:len(result)-w.backupCount]
	}
	return result
}

// moveToBackup renames file to backup name
func (w *TimeFileLogWriter) moveToBackup() error {
	_, err := os.Lstat(w.filename)
	if err == nil { // file exists
		// get the time that this sequence started at and make it a TimeTuple
		t := time.Unix(w.rolloverAt-w.interval, 0).Local()
		fname := w.baseFilename + "." + strftime.Format(w.suffix, t)

		// remove the file with fname if exist
		if _, err := os.Stat(fname); err == nil {
			err = os.Remove(fname)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		}

		// Rename the file to its newfound home
		err = os.Rename(w.baseFilename, fname)
		if err != nil {
			return fmt.Errorf("Rotate: %s\n", err)
		}
	}
	return nil
}

func (w *TimeFileLogWriter) adjustRolloverAt() {
	currTime := time.Now()
	newRolloverAt := w.computeRollover(currTime)

	for newRolloverAt <= currTime.Unix() {
		newRolloverAt = newRolloverAt + w.interval
	}

	w.rolloverAt = newRolloverAt
}

// remove files, according to backupCount
func (w *TimeFileLogWriter) deleteFiles(){
	if w.backupCount > 0 {
		for _, fileName := range w.getFilesToDelete() {
			os.Remove(fileName)
		}
	}
}

// If this is called in a threaded context, it MUST be synchronized
func (w *TimeFileLogWriter) intRotate() error {
	// Close any log file that may be open
	if w.file != nil {
		w.file.Close()
	}

	if w.shouldRollover() {
		// rename file to backup name
		if err := w.moveToBackup(); err != nil {
			return err
		}
	}
	
	// asynchronously delete legacy files to avoid from blocking
	go w.deleteFiles()

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	w.file = fd

	// adjust rolloverAt
	w.adjustRolloverAt()

	return nil
}

// SetFormat sets the logging format (chainable).  Must be called before the first log
// message is written.
func (w *TimeFileLogWriter) SetFormat(format string) *TimeFileLogWriter {
	w.format = format
	return w
}

// Name gets file name
func (w *TimeFileLogWriter) Name() string {
	return w.filename
}

// QueueLen gets rec channel length
func (w *TimeFileLogWriter) QueueLen() int {
	return len(w.rec)
}
