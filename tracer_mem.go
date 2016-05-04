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

// A TracerMemory represents a container for service event entries that stores events
// in-memory.
type TracerMemory struct {
	FilterLevel Level
	Entries     []TracerEntry
}

// NewTracerMemory creates a new instance of TracerMemory type.
func NewTracerMemory(level Level) *TracerMemory {
	return &TracerMemory{
		level,
		make([]TracerEntry, 0),
	}
}

// AddEntry appends a new event to current container.
func (t *TracerMemory) AddEntry(
	level Level, code, msg string,
	htStatus int, err error,
	svcname string, stack ...string,
) *TracerEntry {
	if level < t.FilterLevel {
		return nil
	}

	entry := TracerEntry{
		code,
		level,
		svcname,
		stack,
		msg,
		htStatus,
		err,
	}
	t.Entries = append(t.Entries, entry)
	return &entry
}

// FilterEntries returns event entries where its severity level is greater or equal to
// specified level.
func (t *TracerMemory) FilterEntries(level Level) []TracerEntry {
	var res []TracerEntry
	for _, e := range t.Entries {
		if e.Level >= level {
			res = append(res, e)
		}
	}

	return res
}

var _ Tracer = (*TracerMemory)(nil)
