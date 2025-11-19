// Copyright (C) 2025 The go-concord Authors.
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

package plugins

import (
	"github.com/cybergarage/go-concord/concord/cluster"
	"github.com/cybergarage/go-concord/concord/core"
)

// Coordinator is an interface for the coordinator.
type Coordinator = core.Coordinator

// Observer is an interface for the coordinator observer.
type Observer = core.Observer

// MessageQueue is an interface for the coordinator message queue.
type MessageQueue = core.MessageQueue

// StateType represents the coordinator state type.
type StateType = core.StateType

// Message represents a coordinator message.
type Message = core.Message

// ResultSet represents a coordinator result set.
type ResultSet = core.ResultSet

// Key represents a coordinator key.
type Key = core.Key

// Transaction represents a coordinator transaction.
type Transaction = core.Transaction

// Object represents a coordinator object.
type Object = core.Object

// Service is an interface for the coordinator service.
type Service interface {
	Coordinator
	// SetNodeState posts the specified node status to the
	SetNodeState(node cluster.Node) error
	// GetClusterState gets the current cluster state.
	GetClusterState(cluster string) (cluster.Cluster, error)
}
