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

// A TracerNop represents a container for service event entries that stores nothing.
type TracerNop struct{}

// NewTracerNop creates a new instance of TracerNop type.
func NewTracerNop() *TracerNop {
	return &TracerNop{}
}

// AddEntry does nothing.
func (*TracerNop) AddEntry(Level, string, string, int, error, string, ...string) *TracerEntry {
	return nil
}

var _ Tracer = (*TracerNop)(nil)
