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

package etcd

import (
	"github.com/cybergarage/go-concord/concord/store"
)

// KeyCoder represents a key coder for document keys.
type KeyCoder = store.KeyCoder

type etcdStore struct {
	KeyCoder
}

// Store represents a document store.
type Store = store.Store

// StoreOption represents a function type for etcd store options.
type StoreOption func(*etcdStore)

// WithStoreKeyCoder returns a store option with the specified key coder.
func WithStoreKeyCoder(coder KeyCoder) StoreOption {
	return func(store *etcdStore) {
		store.KeyCoder = coder
	}
}

// NewStore returns a new etcd store instance.
func NewStore(opts ...StoreOption) Store {
	store := &etcdStore{}
	for _, opt := range opts {
		opt(store)
	}
	return store
}

// ServiceName returns the plug-in service name.
func (store *etcdStore) ServiceName() string {
	return "etcd"
}

// Transact begin a new transaction.
func (store *etcdStore) Transact() (store.Transaction, error) {
	return NewTransaction(), nil
}

// Start starts this etcd store.
func (store *etcdStore) Start() error {
	return nil
}

// Stop stops this etcd store.
func (store *etcdStore) Stop() error {
	return nil
}
