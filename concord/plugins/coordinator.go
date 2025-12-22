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

package plugins

import (
	"errors"
	"time"

	"github.com/cybergarage/go-concord/concord/cluster"
	"github.com/cybergarage/go-concord/concord/coordinator"
	"github.com/cybergarage/go-concord/concord/store"
)

type BaseCoordinator struct {
	coordinator.KeyCoder
	cluster.Node
	*time.Ticker
	observers []coordinator.Observer
}

// NewBaseCoordinator returns a new base coordinator instance.
func NewBaseCoordinator() *BaseCoordinator {
	return &BaseCoordinator{
		KeyCoder:  nil,
		Node:      cluster.NewNode(),
		Ticker:    time.NewTicker(time.Second),
		observers: make([]coordinator.Observer, 0),
	}
}

// SetKeyCoder sets the key coder.
func (coord *BaseCoordinator) SetKeyCoder(coder coordinator.KeyCoder) {
	coord.KeyCoder = coder
}

// SetNode sets the coordinator node.
func (coord *BaseCoordinator) SetNode(node cluster.Node) {
	coord.Node = node
}

// AddObserver adds the specified observer.
func (coord *BaseCoordinator) AddObserver(newObserver coordinator.Observer) error {
	for _, observer := range coord.observers {
		if observer == newObserver {
			return nil
		}
	}
	coord.observers = append(coord.observers, newObserver)
	return nil
}

// SetStateObject sets the state object for the specified key.
func (coord *BaseCoordinator) SetStateObject(t coordinator.StateType, obj coordinator.Object) error {
	return errors.New("SetStateObject not implemented")
}

// StateObject gets the object for the specified key and state type.
func (coord *BaseCoordinator) StateObject(t coordinator.StateType, key coordinator.Key) (coordinator.Object, error) {
	return nil, errors.New("StateObject not implemented")
}

// StateObjects gets the result set for the specified key and state type.
func (coord *BaseCoordinator) StateObjects(t coordinator.StateType) (store.ResultSet, error) {
	return nil, errors.New("StateObjects not implemented")
}

// PostMessage posts the specified message to the coordinator.
func (coord *BaseCoordinator) PostMessage(msg coordinator.Message) error {
	return errors.New("PostMessage not implemented")
}

// SetNodeState posts the specified node status to the coordinator.
func (coord *BaseCoordinator) SetNodeState(node cluster.Node) error {
	return errors.New("SetNodeState not implemented")
}

// ClusterState gets the current cluster state.
func (coord *BaseCoordinator) ClusterState(name string) (cluster.Cluster, error) {
	return nil, errors.New("ClusterState not implemented")
}
