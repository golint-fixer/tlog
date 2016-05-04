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

package tlog

// Level represents the severity level of an event.
type Level int

const (
	// LevelTrace defines the severity level for fine-grained debugging events.
	LevelTrace Level = 100

	// LevelDebug defines the severity level for debugging events.
	LevelDebug Level = 200

	// LevelInfo defines the severity level for informational events.
	LevelInfo Level = 300

	// LevelWarn defines the severity level for events that indicates potentially
	// harmful situation.
	LevelWarn Level = 400

	// LevelError defines the severity level for events that indicates failure.
	LevelError Level = 500

	// LevelFatal defines the severity level for very severe errors that leads the
	// service to abort.
	LevelFatal Level = 600
)

var levelText = map[Level]string{
	LevelTrace: "trace",
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warning",
	LevelError: "error",
	LevelFatal: "fatal",
}

// LevelText returns a text for the severity level code. It returns a empty string if the
// code is unknown.
func LevelText(l Level) string {
	return levelText[l]
}
