// Copyright (C) 2022 The Concord Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memdb

import (
	"fmt"

	"github.com/cybergarage/go-concord/concord/coordinator"
	"github.com/cybergarage/go-concord/concord/document"
	"github.com/cybergarage/go-concord/concord/store"
	"github.com/hashicorp/go-memdb"
)

// Memdb represents a Memdb instance.
type resultSet struct {
	coordinator.KeyCoder
	it       memdb.ResultIterator
	key      coordinator.Key
	obj      coordinator.Object
	offset   uint
	limit    uint
	nRead    uint
	lastElem any
}

func newResultSet(coder coordinator.KeyCoder, key coordinator.Key, it memdb.ResultIterator, offset uint, limit uint) store.ResultSet {
	return &resultSet{
		KeyCoder: coder,
		it:       it,
		key:      key,
		obj:      nil,
		offset:   offset,
		limit:    limit,
		nRead:    0,
		lastElem: nil,
	}
}

// Next moves the cursor forward next object from its current position.
func (rs *resultSet) Next() bool {
	if store.NoLimit < rs.limit && uint(rs.limit) <= rs.nRead {
		return false
	}

	for rs.nRead < rs.offset {
		elem := rs.it.Next()
		if elem == nil {
			return false
		}
		rs.nRead++
	}

	elem := rs.it.Next()
	if elem == nil {
		return false
	}
	rs.lastElem = elem
	rs.nRead++

	return true
}

// Object returns an object in the current position.
func (rs *resultSet) Object() (coordinator.Object, error) {
	doc, ok := rs.lastElem.(*Document)
	if !ok {
		return nil, fmt.Errorf("invalid element type: %T", rs.lastElem)
	}
	key, err := rs.DecodeKey([]byte(doc.Key))
	if err != nil {
		return nil, err
	}
	rs.obj = document.NewObjectWith(key, doc.Value)
	return rs.obj, nil
}

// Close closes this result set.
func (rs *resultSet) Close() error {
	// No resources to release in this implementation.
	return nil
}

// Err returns the error, if any, that was encountered during iteration.
// The memdb implementation reads from an in-memory iterator and does not
// accumulate iteration errors, so Err always returns nil.
func (rs *resultSet) Err() error {
	return nil
}
