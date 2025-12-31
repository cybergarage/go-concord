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

import "github.com/cybergarage/go-concord/concord/store"

const (
	createdAtKey      = "_created_at"
	updatedAtKey      = "_updated_at"
	createRevisionKey = "_create_revision"
	modRevisionKey    = "_mod_revision"
	versionKey        = "_version"
)

// storeImpl represents a store implementation.
type storeImpl struct {
	store.Store
}

// NewStoreWith creates a new store with the given object.
func NewStoreWith(store store.Store) Store {
	return &storeImpl{
		Store: store,
	}
}

func (s *storeImpl) Transact() (Transaction, error) {
	txn, err := s.Store.Transact()
	if err != nil {
		return nil, err
	}
	return newTransactionWith(txn), nil
}

// txnImpl represents a transaction implementation.
type txnImpl struct {
	store.Transaction
}

func newTransactionWith(txn store.Transaction) Transaction {
	return &txnImpl{
		Transaction: txn,
	}
}

func (t *txnImpl) Scan(key Key, opts ...StoreOption) (ResultSet, error) {
	resultSet, err := t.Transaction.Scan(key, opts...)
	if err != nil {
		return nil, err
	}
	return newResultSetWith(resultSet), nil
}

// resultSetImpl represents a result set implementation.
type resultSetImpl struct {
	store.ResultSet
}

func newResultSetWith(resultSet store.ResultSet) ResultSet {
	return &resultSetImpl{
		ResultSet: resultSet,
	}
}
