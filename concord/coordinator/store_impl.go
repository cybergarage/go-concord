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

const (
	createdAtKey      = "_created_at"
	updatedAtKey      = "_updated_at"
	createRevisionKey = "_create_revision"
	modRevisionKey    = "_mod_revision"
	versionKey        = "_version"
)

type storeImpl struct {
	Store
}

// NewStoreWith creates a new store with the given object.
func NewStoreWith(store Store) Store {
	return &storeImpl{
		Store: store,
	}
}
