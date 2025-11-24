// Copyright (C) 2022 The go-concord Authors.
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
	_ "embed"
	"errors"
	"fmt"
	"testing"

	"github.com/cybergarage/go-concord/concord"
	"github.com/cybergarage/go-concord/concord/coordinator"
	"github.com/cybergarage/go-concord/concord/document"
	"github.com/cybergarage/go-pict/pict"
	"github.com/cybergarage/go-serix/serix/plugins/coder/key/tuple"
)

//go:embed go_types.pict
var goTypes []byte

func newTestKeyCoder() coordinator.KeyCoder {
	return tuple.NewCoder()
}

func generateCoordinatorObjects() ([]document.Object, error) {
	pict := pict.NewParserWithBytes(goTypes)
	err := pict.Parse()
	if err != nil {
		return []document.Object{}, err
	}

	keys := make([]document.Key, len(pict.Cases()))
	for i, pictCase := range pict.Cases() {
		key := document.NewKey()
		key = append(key, fmt.Sprintf("key%d", i))
		for j, pictParam := range pict.Params() {
			paramType, err := pictParam.Type()
			if err != nil {
				return []document.Object{}, err
			}
			kv, err := pictCase[j].CastTo(paramType)
			if err != nil {
				return []document.Object{}, err
			}
			key = append(key, fmt.Sprintf("%v", kv))
		}
		keys[i] = key
	}

	vals := make([]any, len(pict.Cases()))
	for i, pictCase := range pict.Cases() {
		val := map[string]any{}
		for j, pictParam := range pict.Params() {
			pictElem := pictCase[j]
			paramName := pictParam.Name()
			paramType, err := pictParam.Type()
			if err != nil {
				return []document.Object{}, err
			}
			v, err := pictElem.CastTo(paramType)
			if err != nil {
				return []coordinator.Object{}, err
			}
			val[paramName] = v
		}
		vals[i] = val
	}

	objs := make([]coordinator.Object, len(pict.Cases()))
	for n, key := range keys {
		obj, err := document.NewObjectFrom(key, vals[n])
		if err != nil {
			return []coordinator.Object{}, err
		}
		objs[n] = obj
	}

	return objs, nil
}

func updateCoordinatorObjects(objs []coordinator.Object) ([]coordinator.Object, error) {
	pict := pict.NewParserWithBytes(goTypes)
	err := pict.Parse()
	if err != nil {
		return []coordinator.Object{}, err
	}

	for i, pictCase := range pict.Cases() {
		val := []any{}
		for j, pictParam := range pict.Params() {
			pictElem := pictCase[j]
			paramType, err := pictParam.Type()
			if err != nil {
				return []document.Object{}, err
			}
			v, err := pictElem.CastTo(paramType)
			if err != nil {
				return []coordinator.Object{}, err
			}
			val = append(val, v)
		}
		obj, err := document.NewObjectFrom(objs[i].Key(), val)
		if err != nil {
			return []coordinator.Object{}, err
		}
		objs[i] = obj
	}

	return objs, nil
}

// CoordinatorStoreTest tests the coordinator store functionality.
// nolint:goerr113, gocognit, gci, gocyclo, gosec, maintidx
func CoordinatorStoreTest(t *testing.T, coord concord.Service) {
	t.Helper()

	coord.SetKeyCoder(newTestKeyCoder())

	cancel := func(t *testing.T, txn coordinator.Transaction) {
		t.Helper()
		if err := txn.Cancel(); err != nil {
			t.Error(err)
		}
	}

	// Starts the coordinator service

	if err := coord.Start(); err != nil {
		t.Error(err)
		return
	}

	// Terminates the coordinator service

	defer func() {
		if err := coord.Stop(); err != nil {
			t.Error(err)
		}
	}()

	// Generates test keys and objects

	objs, err := generateCoordinatorObjects()
	if err != nil {
		t.Error(err)
		return
	}

	// Inserts new objects

	for _, obj := range objs {
		tx, err := coord.Transact()
		if err != nil {
			t.Error(err)
			return
		}
		if err := tx.Set(obj); err != nil {
			cancel(t, tx)
			t.Error(err)
			return
		}
		if err := tx.Commit(); err != nil {
			t.Error(err)
			return
		}
	}

	// Selects inserted objects

	for _, obj := range objs {
		tx, err := coord.Transact()
		if err != nil {
			t.Error(err)
			return
		}
		retObj, err := tx.Get(obj.Key())
		if err != nil {
			cancel(t, tx)
			t.Error(err)
			return
		}
		if !retObj.Equals(obj) {
			cancel(t, tx)
			t.Error(err)
			return
		}
		if err := tx.Commit(); err != nil {
			t.Error(err)
			return
		}
	}

	// Updates inserted objects

	objs, err = updateCoordinatorObjects(objs)
	if err != nil {
		t.Error(err)
		return
	}

	for _, obj := range objs {
		tx, err := coord.Transact()
		if err != nil {
			t.Error(err)
			return
		}
		if err := tx.Set(obj); err != nil {
			cancel(t, tx)
			t.Error(err)
			return
		}
		if err := tx.Commit(); err != nil {
			t.Error(err)
			return
		}
	}

	// Selects update objects

	for _, obj := range objs {
		tx, err := coord.Transact()
		if err != nil {
			t.Error(err)
			return
		}
		retObj, err := tx.Get(obj.Key())
		if err != nil {
			cancel(t, tx)
			t.Error(err)
			return
		}

		if !retObj.Equals(obj) {
			cancel(t, tx)
			t.Error(err)
			return
		}
		if err := tx.Commit(); err != nil {
			t.Error(err)
			return
		}
	}

	// Removes updated objects

	for _, obj := range objs {
		tx, err := coord.Transact()
		if err != nil {
			t.Error(err)
			return
		}
		err = tx.Remove(obj.Key())
		if err != nil {
			cancel(t, tx)
			t.Error(err)
			return
		}
		_, err = tx.Get(obj.Key())
		if !errors.Is(err, coordinator.ErrNotExist) {
			cancel(t, tx)
			t.Error(err)
			return
		}
		if err := tx.Commit(); err != nil {
			t.Error(err)
			return
		}
	}

	// Truncates all objects

	tx, err := coord.Transact()
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.Truncate()
	if err != nil {
		cancel(t, tx)
		t.Error(err)
		return
	}
	if err := tx.Commit(); err != nil {
		t.Error(err)
		return
	}
}
