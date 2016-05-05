/*
 * Copyright 2016 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package services

import (
	"fmt"
	"time"

	"github.com/raiqub/tlog"
	"gopkg.in/mgo.v2"
)

const (
	// ErrInsertLogEntryCode defines the log entry code for database errors.
	ErrInsertLogEntryCode = "insert_db_log_entry"

	// ErrInsertLogEntryMessage defines the log entry message for database errors.
	ErrInsertLogEntryMessage = "Error inserting a new log entry"

	// LoggerServiceName defines the service name for events from LoggerMongo.
	LoggerServiceName = "LoggerMongo"
)

// A LoggerMongo represents a service that allows to log events to MongoDB database.
type LoggerMongo struct {
	appName    string
	appVersion string
	col        *mgo.Collection
}

// NewLoggerMongo creates a new instance of LoggerMongo.
func NewLoggerMongo(col *mgo.Collection) *LoggerMongo {
	return &LoggerMongo{"", "", col}
}

// LogEntry adds a new logging entry.
func (l *LoggerMongo) LogEntry(entry *Entry) {
	fmt.Println(entry.ToLogfmt())

	if err := l.col.Insert(entry); err != nil {
		l.logDBError(err)
	}
}

// NewContext creates a new context that allows to services log its events not
// knowing current context.
func (l *LoggerMongo) NewContext() *TracerContext {
	return NewEntryContext(l, l.appName, l.appVersion)
}

// SetApplication sets the name and version of current application.
func (l *LoggerMongo) SetApplication(name, version string) {
	l.appName = name
	l.appVersion = version
}

func (l *LoggerMongo) logDBError(err error) {
	entry := Entry{
		Level:      tlog.LevelError,
		Time:       time.Now(),
		AppName:    l.appName,
		AppVersion: l.appVersion,
		SvcName:    LoggerServiceName,
		Code:       ErrInsertLogEntryCode,
		Message:    ErrInsertLogEntryMessage,
	}
	entry.SetError(err)
	fmt.Println(entry.ToLogfmt())
}

var _ Logger = (*LoggerMongo)(nil)
