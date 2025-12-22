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

package concord

import (
	"github.com/cybergarage/go-concord/concord/cluster"
	"github.com/cybergarage/go-concord/concord/coordinator"
)

// Service represents a coordinator service interface.
type Service interface {
	// ServiceName returns the name of the coordinator service.
	ServiceName() string
	// Coordinator represents the coordinator service interface.
	coordinator.Coordinator
	// Start starts the coordinator service.
	Start() error
	// Stop stops the coordinator service.
	Stop() error
	// SetNodeState posts the specified node status to the coordinator.
	SetNodeState(node cluster.Node) error
	// ClusterState gets the current cluster state.
	ClusterState(cluster string) (cluster.Cluster, error)
}

// ServiceFactory is a function type for creating a new coordinator service.
type ServiceFactory func() Service
