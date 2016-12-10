// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package csdb

import (
	"sync"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/errors"
)

// @deprecated
const (
	MainTable       = "main_table"
	AdditionalTable = "additional_table"
	ScopeTable      = "scope_table"
	ghostTableName  = "PendingTableName"
)

// TableOption applies options to the Tables struct.
type TableOption func(*Tables) error

// Tables handles all the tables defined for a package. Thread safe.
type Tables struct {
	// Schema represents the name of the database. Might be empty.
	Schema string
	// Prefix will be put in front of each table name. TODO implement table prefix.
	Prefix string
	mu     sync.RWMutex
	// ts uses int as the table index.
	// What is the reason to use int as the table index and not a name? Because
	// table names between M1 and M2 get renamed and in a Go SQL code generator
	// script of the CoreStore project, we can guarantee that the generated
	// index constant will always stay the same but the name of the table
	// differs.
	ts map[int]*Table
}

// WithTable inserts a new table to the Tables struct, identified by its index.
// You can optionally specify the columns. What is the reason to use int as the
// table index and not a name? Because table names between M1 and M2 get renamed
// and in a Go SQL code generator script of the CoreStore project, we can
// guarantee that the generated index constant will always stay the same but the
// name of the table differs.
func WithTable(idx int, tableName string, cols ...*Column) TableOption {
	return func(tm *Tables) error {
		if err := IsValidIdentifier(tableName); err != nil {
			return errors.Wrap(err, "[csdb] WithNewTable.IsValidIdentifier")
		}

		if err := tm.Upsert(idx, NewTable(tableName, cols...)); err != nil {
			return errors.Wrap(err, "[csdb] WithNewTable.Tables.Insert")
		}
		return nil
	}
}

// WithTableLoadColumns inserts a new table to the Tables struct, identified by
// its index. What is the reason to use int as the table index and not a name?
// Because table names between M1 and M2 get renamed and in a Go SQL code
// generator script of the CoreStore project, we can guarantee that the
// generated index constant will always stay the same but the name of the table
// differs.
func WithTableLoadColumns(db dbr.Querier, idx int, tableName string) TableOption {
	return func(tm *Tables) error {
		if err := IsValidIdentifier(tableName); err != nil {
			return errors.Wrap(err, "[csdb] WithTableLoadColumns.IsValidIdentifier")
		}

		t := NewTable(tableName)
		t.Schema = tm.Schema
		if err := t.LoadColumns(db); err != nil {
			return errors.Wrap(err, "[csdb] WithTableLoadColumns.LoadColumns")
		}

		if err := tm.Upsert(idx, t); err != nil {
			return errors.Wrap(err, "[csdb] Tables.Insert")
		}
		return nil
	}
}

// WithTableNames creates for each table name and its index a new table pointer.
// You should call afterwards the functional option WithLoadColumnDefinitions.
// This function returns an error if a table index already exists.
func WithTableNames(idx []int, tableName []string) TableOption {
	return func(tm *Tables) error {
		if len(idx) != len(tableName) {
			return errors.NewNotValidf("[csdb] Length of the index must be equal to the length of the table names: %d != %d", len(idx), len(tableName))
		}

		if err := IsValidIdentifier(tableName...); err != nil {
			return errors.Wrap(err, "[csdb] WithTable.IsValidIdentifier")
		}

		for i, tn := range tableName {
			if err := tm.Upsert(idx[i], NewTable(tn)); err != nil {
				return errors.Wrapf(err, "[csdb] Tables.Insert %q", tn)
			}
		}
		return nil
	}
}

// WithLoadColumnDefinitions loads the column definitions from the database for each
// table in the internal map. Thread safe.
func WithLoadColumnDefinitions(db dbr.Querier) TableOption {
	return func(tm *Tables) error {
		tm.mu.Lock()
		defer tm.mu.Unlock()

		for _, table := range tm.ts {
			// TODO(CyS) could be refactored to fire only one query ... but later.
			if err := table.LoadColumns(db); err != nil {
				return errors.Wrap(err, "[csdb] table.LoadColumns")
			}
		}
		return nil
	}
}

// WithTableDMLListeners adds event listeners to a table object. It doesn't
// matter if the table has already been set. If the table object gets set later,
// the events will be copied to the new object.
func WithTableDMLListeners(idx int, events ...*dbr.ListenerBucket) TableOption {
	return func(tm *Tables) error {
		tm.mu.Lock()
		defer tm.mu.Unlock()

		t, ok := tm.ts[idx]
		if !ok || t == nil {
			t = NewTable(ghostTableName)
		}
		t.ListenerBucket.Merge(events...)
		tm.ts[idx] = t

		return nil
	}
}

// NewTables creates a new TableService satisfying interface Manager.
func NewTables(opts ...TableOption) (*Tables, error) {
	tm := &Tables{
		ts: make(map[int]*Table),
	}
	if err := tm.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[csdb] NewTables applied option error")
	}
	return tm, nil
}

// MustNewTables same as NewTableService but panics on error.
func MustNewTables(opts ...TableOption) *Tables {
	ts, err := NewTables(opts...)
	if err != nil {
		panic(err)
	}
	return ts
}

// MustInitTables helper function in init() statements to initialize the global
// table collection variable independent of knowing when this variable is nil.
// We cannot assume the correct order, how all init() invocations are executed,
// at least they don't run in parallel during packet initialization. Yes ... bad
// practice to rely on init ... but for now it works very well.
// TODO(CyS) rethink and refactor maybe.
func MustInitTables(ts *Tables, opts ...TableOption) *Tables {
	if ts == nil {
		var err error
		ts, err = NewTables()
		if err != nil {
			panic(err)
		}
	}
	if err := ts.Options(opts...); err != nil {
		panic(err)
	}
	return ts
}

// Options applies options to the Tables service.
func (tm *Tables) Options(opts ...TableOption) error {
	for _, o := range opts {
		if err := o(tm); err != nil {
			return errors.Wrap(err, "[csdb] Applied option error")
		}
	}
	return nil
}

// Table returns the structure from a map m by a giving index i. What is the
// reason to use int as the table index and not a name? Because table names
// between M1 and M2 get renamed and in a Go SQL code generator script of the
// CoreStore project, we can guarantee that the generated index constant will
// always stay the same but the name of the table differs.
func (tm *Tables) Table(i int) (*Table, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if t, ok := tm.ts[i]; ok {
		return t, nil
	}
	return nil, errors.NewNotFoundf("[csdb] Table at index %d not found.", i)
}

// MustTable same as Table function but panics when the table cannot be found or
// any other error occurs.
func (tm *Tables) MustTable(i int) *Table {
	t, err := tm.Table(i)
	if err != nil {
		panic(err)
	}
	return t
}

// Name is a short hand to return a table name by given index i. Does not return
// an error when the table can't be found but returns an empty string.
func (tm *Tables) Name(i int) string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if ts, ok := tm.ts[i]; ok && ts != nil {
		return ts.Name
	}
	return ""
}

// Len returns the number of all tables.
func (tm *Tables) Len() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.ts)
}

// Upsert adds or updates a new table into the internal cache. *Table cannot be nil.
func (tm *Tables) Upsert(i int, tNew *Table) error {
	_ = tNew.Name // let it panic as early as possible if *Table is nil

	tm.mu.Lock()
	defer tm.mu.Unlock()

	tOld, ok := tm.ts[i]
	if tOld == nil || !ok {
		tm.ts[i] = tNew
		return nil
	}

	if tOld.Name == ghostTableName {
		// for now copy only the events from the existing table
		tNew.ListenerBucket.Merge(&tOld.ListenerBucket)
	}

	tm.ts[i] = tNew
	return nil
}

// Delete removes tables by their given indexes. If no index has been passed
// then all entries get removed and the map reinitialized.
func (tm *Tables) Delete(idxs ...int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for _, idx := range idxs {
		delete(tm.ts, idx)
	}
	if len(idxs) == 0 {
		tm.ts = make(map[int]*Table)
	}
}