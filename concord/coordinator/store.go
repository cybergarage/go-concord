// Copyright (C) 2025 The Concord Authors.
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

package coordinator

import (
	"github.com/cybergarage/go-concord/concord/store"
)

// Key represents a document key.
type Key = store.Key

// KeyCoder represents a key coder for document keys.
type KeyCoder = store.KeyCoder

// Object represents a document object.
type Object = store.Object

// StoreOption represents a store option.
type StoreOption = store.Option

// Result represents a result which includes a single object.
type Result interface {
	Object() (Object, error)
}

// ResultSet represents a result set which includes range operation results.
type ResultSet interface {
	// Next moves the cursor forward next object from its current position.
	Next() bool
	// Result returns the current result object.
	Result() (Result, error)
}

// Transaction represents a transaction interface.
type Transaction interface {
	// Set sets the object for the specified key.
	Set(obj Object) error
	// Get gets the object for the specified key.
	Get(key Key) (Object, error)
	// Scan returns the result set for the specified key.
	Scan(key Key, opts ...StoreOption) (ResultSet, error)
	// Remove removes the object for the specified key.
	Remove(key Key) error
	// Commit commits this transaction.
	Commit() error
	// Cancel cancels this transaction.
	Cancel() error
	// Truncate removes all objects.
	Truncate() error
}

// Store represents a coordination store inteface.
type Store interface {
	// Transact begin a new transaction.
	Transact() (Transaction, error)
	// Start starts this store.
	Start() error
	// Stop stops this store.
	Stop() error
}
