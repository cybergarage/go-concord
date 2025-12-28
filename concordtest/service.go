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

	"github.com/cybergarage/go-concord/concord"
	"github.com/cybergarage/go-concord/concord/store"
)

// Service represents a concord service.
type Service struct {
	concord.Service
}

// NewServiceWith creates a new service with the given object.
func NewServiceWith(service concord.Service) *Service {
	return &Service{
		Service: service,
	}
}

// Store returns the store instance.
func (service *Service) Store() (*Store, error) {
	store, ok := service.Service.(store.Store)
	if !ok {
		return nil, fmt.Errorf("service is not a store")
	}
	return NewStoreWith(store), nil
}
