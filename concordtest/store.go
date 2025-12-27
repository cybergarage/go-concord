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
	"strings"

	"github.com/cybergarage/go-concord/concord/document"
	"github.com/cybergarage/go-concord/concord/store"
)

// Store represents a document store.
type Store struct {
	store.Store
}

// NewStoreWith creates a new store with the given object.
func NewStoreWith(store store.Store) *Store {
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
		_, err := tx.Scan(key)
		if err != nil {
			continue
		}
		/*
			for rs.Next() {
				obj, err := rs.Object()
				if err != nil {
					continue
				}
				keys := obj.Key().Elements()
				keyHeaderBytes, ok := keys[0].([]byte)
				if !ok {
					lines = append(lines, fmt.Sprintf("%v: %v", keys[1:], obj.Value()))
				}
				keyHeader := dockv.NewKeyHeaderFrom(keyHeaderBytes)

				switch keyHeader.Type() {
				case dockv.DatabaseObject, dockv.CollectionObject, dockv.DocumentObject:
					r := bytes.NewReader(obj.Value())
					val, err := docStore.DecodeDocument(r)
					if err != nil {
						lines = append(lines, fmt.Sprintf("%v %v: %v", keyHeader, keys[1:], obj.Value()))
						continue
					}
					lines = append(lines, fmt.Sprintf("%v %v: %v", keyHeader, keys[1:], val))
				case dockv.IndexObject:
					lines = append(lines, fmt.Sprintf("%v %v:", keyHeader, keys[1:]))
				}
			}
		*/
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
