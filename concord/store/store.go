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

package store

import (
	"github.com/cybergarage/go-concord/concord/document"
)

// Key represents a document key.
type Key = document.Key

// KeyCoder represents a key coder for document keys.
type KeyCoder = document.KeyCoder

// Object represents a document object.
type Object = document.Object

// ResultSet represents a result set which includes range operation results.
type ResultSet interface {
	// Next moves the cursor forward next object from its current position.
	Next() bool
	// Object returns an object in the current position.
	Object() Object
}

// Transaction represents a transaction interface.
type Transaction interface {
	// Set sets the object for the specified key.
	Set(obj Object) error
	// Get gets the object for the specified key.
	Get(key Key) (Object, error)
	// Scan returns the result set for the specified key.
	Scan(key Key, opts ...Option) (ResultSet, error)
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
	// SetKeyCoder sets the key coder.
	SetKeyCoder(coder KeyCoder)
	// KeyCoder represents a key decoder and encoder interface.
	KeyCoder
	// Transact begin a new transaction.
	Transact() (Transaction, error)
	// Start starts this store.
	Start() error
	// Stop stops this store.
	Stop() error
}
