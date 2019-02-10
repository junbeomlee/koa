/*
 * Copyright 2018 De-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package translate

import (
	"fmt"
)

type EntryError struct {
	Id string
}

func (e EntryError) Error() string {
	return fmt.Sprintf("[%s] definition doesn't exist", e.Id)
}

type MemTracer interface {
	MemDefiner
	MemEntryGetter
}

// Define() saves an variable to EntryMap and increase the MemoryCounter.
// This should be used when compiles the assign statement.
// ex)
// a = 5 -> Define("a", 5)
// b = "abc" -> Define("b", "abc")
type MemDefiner interface {
	Define(id string, value []byte) MemEntry
}

// MemEntryGetter gets the data of the memory entry.
// GetOffsetOfEntry() returns the offset of the memory entry corresponding the Id.
// GetSizeOfEntry() returns the size of the memory entry corresponding the Id.
type MemEntryGetter interface {
	GetOffsetOfEntry(id string) (int, error)
	GetSizeOfEntry(id string) (int, error)
}

// MemEntry saves size and offset of the value which the variable has.
type MemEntry struct {
	Offset int
	Size   int
}

// MemEntryTable is used to know the location of the memory
type MemEntryTable struct {
	EntryMap      map[string]MemEntry
	MemoryCounter int
}

func NewMemEntryTable() *MemEntryTable {
	return &MemEntryTable{
		EntryMap:      make(map[string]MemEntry),
		MemoryCounter: 0,
	}
}

func (m *MemEntryTable) Define(id string, value []byte) MemEntry {
	entry := MemEntry{
		Offset: m.MemoryCounter,
	}

	size := len(value)
	entry.Size = size
	m.MemoryCounter += size
	m.EntryMap[id] = entry

	return entry
}

func (m MemEntryTable) GetOffsetOfEntry(id string) (int, error) {
	entry, ok := m.EntryMap[id]
	if !ok {
		return 0, EntryError{
			Id: id,
		}
	}

	return entry.Offset, nil
}

func (m MemEntryTable) GetSizeOfEntry(id string) (int, error) {
	entry, ok := m.EntryMap[id]
	if !ok {
		return 0, EntryError{
			Id: id,
		}
	}

	return entry.Size, nil
}
