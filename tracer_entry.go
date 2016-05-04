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

//go:generate ffjson $GOFILE

package tlog

import "bytes"

// A TracerEntry represents a event created by a service.
type TracerEntry struct {
	Code       string   `bson:"code" json:"code"`
	Level      Level    `bson:"level" json:"-"`
	Service    string   `bson:"service" json:"-"`
	Stack      []string `bson:"stack" json:"-"`
	Message    string   `bson:"msg" json:"message"`
	HTTPStatus int      `bson:"-" json:"status"`
	InnerError error    `bson:"inner" json:"-"`
}

// Error returns string representation of current instance error.
func (e *TracerEntry) Error() string {
	var buf bytes.Buffer
	if len(e.Code) > 0 {
		buf.WriteString(e.Code)
		buf.WriteString(": ")
	}
	buf.WriteString(e.Message)

	return buf.String()
}

// String returns string representation of current instance.
func (e *TracerEntry) String() string {
	return e.Error()
}
