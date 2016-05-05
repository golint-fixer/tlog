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
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/raiqub/tlog"
	"gopkg.in/mgo.v2/bson"
)

// A Entry represents a logging entry.
type Entry struct {
	ID           bson.ObjectId `bson:"_id"`
	Level        tlog.Level    `bson:"level"`
	Time         time.Time     `bson:"time"`
	AppName      string        `bson:"app"`
	AppVersion   string        `bson:"appver"`
	SvcName      string        `bson:"svc,omitempty"`
	Stack        []string      `bson:"stack,omitempty"`
	NumCPU       int           `bson:"numcpu,omitempty"`
	NumThreads   int           `bson:"numthr,omitempty"`
	Host         string        `bson:"host,omitempty"`
	ReqMethod    string        `bson:"method,omitempty"`
	ReqURL       string        `bson:"url,omitempty"`
	ReqIP        string        `bson:"ip,omitempty"`
	ReqUserAgent string        `bson:"usragent,omitempty"`
	ReqReferer   string        `bson:"reqref,omitempty"`
	ReqClientIP  string        `bson:"cliip,omitempty"`
	AccountID    string        `bson:"actid,omitempty"`
	AccountName  string        `bson:"actname,omitempty"`
	AccountPlan  string        `bson:"actplan,omitempty"`
	AccountEmail string        `bson:"actmail,omitempty"`
	Code         string        `bson:"code"`
	Message      string        `bson:"msg"`
	InnerError   string        `bson:"err,omitempty"`
}

// SetError decodes specified error to InnerError field.
func (e *Entry) SetError(err error) {
	if err == nil {
		return
	}

	errType := reflect.TypeOf(err)
	var typeName string
	if errType.Kind() == reflect.Ptr {
		typeName = errType.Elem().Name()
	} else {
		typeName = errType.Name()
	}

	e.InnerError = fmt.Sprintf("%s: %s", typeName, err.Error())
}

// ToLogfmt returns the Logfmt encoding of current entry.
func (e *Entry) ToLogfmt() string {
	var buf bytes.Buffer
	buf.WriteString(`time="` + e.Time.Format(time.RFC3339) + `"`)
	buf.WriteString(` level=` + tlog.LevelText(e.Level))
	buf.WriteString(` appname="` + e.AppName + `"`)
	buf.WriteString(` appver=` + e.AppVersion)
	if len(e.SvcName) > 0 {
		buf.WriteString(` svcname=` + e.SvcName)
	}
	if len(e.Stack) > 0 {
		buf.WriteString(` stack="` + strings.Join(e.Stack, ":") + `"`)
	}

	if e.NumCPU > 0 {
		buf.WriteString(` numcpu=` + strconv.Itoa(e.NumCPU))
	}
	if e.NumThreads > 0 {
		buf.WriteString(` numthr=` + strconv.Itoa(e.NumThreads))
	}

	if len(e.Host) > 0 {
		buf.WriteString(` host=` + e.Host)
	}
	if len(e.ReqMethod) > 0 {
		buf.WriteString(` method=` + e.ReqMethod)
		buf.WriteString(` url=` + e.ReqURL)
		buf.WriteString(` ip="` + e.ReqIP + `"`)
		buf.WriteString(` usragent="` + e.ReqUserAgent + `"`)
		buf.WriteString(` reqref="` + e.ReqReferer + `"`)
		buf.WriteString(` cliip="` + e.ReqClientIP + `"`)
	}

	if len(e.AccountID) > 0 {
		buf.WriteString(` actid="` + e.AccountID + `"`)
		buf.WriteString(` actname="` + e.AccountName + `"`)
		buf.WriteString(` actplan="` + e.AccountPlan + `"`)
		buf.WriteString(` actmail="` + e.AccountEmail + `"`)
	}

	buf.WriteString(` code=` + e.Code)
	buf.WriteString(` msg="` + e.Message + `"`)
	if len(e.InnerError) > 0 {
		buf.WriteString(` err="` + e.InnerError + `"`)
	}
	return buf.String()
}
