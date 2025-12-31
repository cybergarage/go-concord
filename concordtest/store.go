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

package concordtest

import (
	"fmt"
	"strings"

	"github.com/cybergarage/go-concord/concord/coordinator"
	"github.com/cybergarage/go-concord/concord/document"
)

// Store represents a document store.
type Store struct {
	coordinator.Store
}

// NewStoreWith creates a new store with the given object.
func NewStoreWith(store coordinator.Store) *Store {
	return &Store{
		Store: store,
	}
}

// Dump returns a string representation of the store.
func (s *Store) Dump() ([]string, error) {
	lines := []string{}

	tx, err := s.Store.Transact()
	if err != nil {
		return lines, err
	}

	keys := document.HeaderPrefixKeys()

	for _, key := range keys {
		rs, err := tx.Scan(key)
		if err != nil {
			continue
		}
		for rs.Next() {
			r, err := rs.Result()
			if err != nil {
				continue
			}
			obj, err := r.Object()
			if err != nil {
				continue
			}
			lines = append(lines, fmt.Sprintf("%v: %v", obj.Key(), obj))
		}
		if err := rs.Close(); err != nil {
			continue
		}
	}

	return lines, tx.Commit()
}

// String returns a string representation of the store.
func (s *Store) String() string {
	lines, err := s.Dump()
	if err != nil {
		return ""
	}
	return strings.Join(lines, "\n")
}
