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
	"testing"

	"github.com/cybergarage/go-concord/concord"
	"github.com/cybergarage/go-concord/concord/plugins/memdb"
	"github.com/cybergarage/go-logger/log"
)

func TestServices(t *testing.T) {
	log.SetSharedLogger(log.NewStdoutLogger(log.LevelInfo))

	factories := []concord.ServiceFactory{
		memdb.NewService,
	}

	for _, factory := range factories {
		services := []*Service{}
		for i := range 2 {
			service, err := factory()
			if err != nil {
				t.Error(err)
				return
			}
			service.SetHost(fmt.Sprintf("service%02d", i+1))
			if err := service.Start(); err != nil {
				t.Error(err)
				return
			}
			services = append(services, NewServiceWith(service))
		}

		ValidateServiceCoordinationSuite(t, services)
	}
}

// ValidateServiceCoordinationSuite validates the service coordination functionality.
func ValidateServiceCoordinationSuite(t *testing.T, services []*Service) {
	t.Helper()
	name := services[0].ServiceName()
	t.Run(name, func(t *testing.T) {
		tests := []struct {
			name string
			fn   func(t *testing.T, coords []*Service)
		}{
			{"messaging", ValidateCoordinatorMessageFlow},
			{"clustring", ValidateCoordinatorClusterState},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				test.fn(t, services)
			})
		}
	})

	for _, coord := range services {
		if err := coord.Stop(); err != nil {
			t.Error(err)
			return
		}
	}
}
