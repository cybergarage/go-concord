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
	"errors"

	"github.com/cybergarage/go-concord/concord/store"
	"github.com/hashicorp/go-memdb"
)

const (
	tableName   = "coordinator"
	idName      = "id"
	idFieldName = "Key"
	prefix      = "_prefix"
)

// KeyCoder represents a key coder for document keys.
type KeyCoder = store.KeyCoder

// Store represents a document store.
type Store = store.Store

var sharedMemDB *memdb.MemDB = nil

type memdbStore struct {
	keyCorder KeyCoder
	*memdb.MemDB
}

// StoreOption represents a function type for memdb store options.
type StoreOption func(*memdbStore) error

// WithStoreKeyCoder returns a store option with the specified key coder.
func WithStoreKeyCoder(coder KeyCoder) StoreOption {
	return func(store *memdbStore) error {
		store.keyCorder = coder
		return nil
	}
}

// NewStore returns a new memdb store instance.
func NewStore(opts ...StoreOption) (Store, error) {
	store := &memdbStore{
		keyCorder: nil,
		MemDB:     nil,
	}
	for _, opt := range opts {
		if err := opt(store); err != nil {
			return nil, err
		}
	}
	if store.keyCorder == nil {
		return nil, errors.New("key coder is not set")
	}
	return store, nil
}

func (store *memdbStore) Transact() (store.Transaction, error) {
	return newTransactionWith(store.keyCorder, store.MemDB.Txn(true)), nil
}

// Start starts this memdb store.
func (store *memdbStore) Start() error {
	if sharedMemDB != nil {
		store.MemDB = sharedMemDB
		return nil
	}

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableName: {
				Name: tableName,
				Indexes: map[string]*memdb.IndexSchema{
					idName: {
						Name:         idName,
						AllowMissing: false,
						Unique:       true,
						Indexer:      &BinaryFieldIndexer{},
					},
				},
			},
		},
	}
	memDB, err := memdb.NewMemDB(schema)
	if err != nil {
		return errors.Join(err, store.Stop())
	}
	sharedMemDB = memDB
	store.MemDB = sharedMemDB
	return nil
}

// Stop stops this etcd coordinator.
func (store *memdbStore) Stop() error {
	return nil
}
