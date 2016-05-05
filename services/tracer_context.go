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
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/raiqub/tlog"
	"gopkg.in/mgo.v2/bson"
)

// A TracerContext represents a Tracer context that allows to services log its events not
// knowing current context.
type TracerContext struct {
	logger     Logger
	sync       bool
	appName    string
	appVersion string
	filter     tlog.Level

	request *http.Request

	actID    string
	actName  string
	actPlan  string
	actEmail string
}

// NewEntryContext creates a new instance of TracerContext.
func NewEntryContext(logger Logger, appName, appVersion string) *TracerContext {
	return &TracerContext{
		logger, false, appName, appVersion, 0,
		nil,
		"", "", "", "",
	}
}

// AddEntry appends a new event to current container.
func (ctx *TracerContext) AddEntry(
	level tlog.Level, code, msg string,
	htStatus int, err error,
	svcname string, stack ...string,
) *tlog.TracerEntry {
	if level < ctx.filter {
		return nil
	}

	tentry := tlog.TracerEntry{
		Code:       code,
		Level:      level,
		Service:    svcname,
		Stack:      stack,
		Message:    msg,
		HTTPStatus: htStatus,
		InnerError: err,
	}

	if ctx.sync {
		ctx.logEntry(tentry)
	} else {
		go ctx.logEntry(tentry)
	}

	return &tentry
}

// SetAccount sets user account details.
func (ctx *TracerContext) SetAccount(
	actID string,
	actName string,
	actPlan string,
	actEmail string,
) {
	ctx.actID = actID
	ctx.actName = actName
	ctx.actPlan = actPlan
	ctx.actEmail = actEmail
}

// SetFilter sets the minimum severity level for logging. Events that its severity level
// is less than specified value will be ignored.
func (ctx *TracerContext) SetFilter(level tlog.Level) {
	ctx.filter = level
}

// SetRequest sets the HTTP request to get context values.
func (ctx *TracerContext) SetRequest(req *http.Request) {
	ctx.request = req
}

// SetSync determines whether log writes are synchronous or not. Default false
// (asynchronous).
func (ctx *TracerContext) SetSync(status bool) {
	ctx.sync = status
}

func (ctx *TracerContext) logEntry(tentry tlog.TracerEntry) {
	entry := Entry{
		ID:         bson.NewObjectId(),
		Time:       time.Now(),
		AppName:    ctx.appName,
		AppVersion: ctx.appVersion,
		NumCPU:     runtime.NumCPU(),
		NumThreads: runtime.NumGoroutine(),

		Level:   tentry.Level,
		SvcName: tentry.Service,
		Stack:   tentry.Stack,
		Code:    tentry.Code,
		Message: tentry.Message,

		AccountID:    ctx.actID,
		AccountName:  ctx.actName,
		AccountPlan:  ctx.actEmail,
		AccountEmail: ctx.actEmail,
	}

	entry.SetError(tentry.InnerError)

	if ctx.request != nil {
		entry.Host = ctx.request.Host
		entry.ReqMethod = ctx.request.Method
		entry.ReqURL = ctx.request.RequestURI
		entry.ReqIP = ctx.request.RemoteAddr
		entry.ReqUserAgent = ctx.request.Header.Get("User-Agent")
		entry.ReqReferer = ctx.request.Header.Get("Referer")
		entry.ReqClientIP = clientIP(ctx.request)
	}

	ctx.logger.LogEntry(&entry)
}

// ClientIP implements a best effort algorithm to return the real client IP, it
// parses X-Real-IP and X-Forwarded-For in order to work properly with
// reverse-proxies such us: nginx or haproxy.
func clientIP(req *http.Request) string {
	clientIP := strings.TrimSpace(req.Header.Get("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = req.Header.Get("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(
		strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

var _ tlog.Tracer = (*TracerContext)(nil)
