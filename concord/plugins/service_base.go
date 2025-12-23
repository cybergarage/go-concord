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

// ServiceBase is a base coordinator service implementation.
type ServiceBase struct {
	coordinator.KeyCoder
	cluster.Node
	*time.Ticker
	observers []coordinator.Observer
}

// NewServiceBase returns a new base coordinator instance.
func NewServiceBase() *ServiceBase {
	return &ServiceBase{
		KeyCoder:  nil,
		Node:      cluster.NewNode(),
		Ticker:    time.NewTicker(time.Second),
		observers: make([]coordinator.Observer, 0),
	}
}

// SetKeyCoder sets the key coder.
func (coord *ServiceBase) SetKeyCoder(coder coordinator.KeyCoder) {
	coord.KeyCoder = coder
}

// SetNode sets the coordinator node.
func (coord *ServiceBase) SetNode(node cluster.Node) {
	coord.Node = node
}

// AddObserver adds the specified observer.
func (coord *ServiceBase) AddObserver(newObserver coordinator.Observer) error {
	for _, observer := range coord.observers {
		if observer == newObserver {
			return nil
		}
	}
	coord.observers = append(coord.observers, newObserver)
	return nil
}

// SetStateObject sets the state object for the specified key.
func (coord *ServiceBase) SetStateObject(t coordinator.StateType, obj coordinator.Object) error {
	return errors.New("SetStateObject not implemented")
}

// StateObject gets the object for the specified key and state type.
func (coord *ServiceBase) StateObject(t coordinator.StateType, key coordinator.Key) (coordinator.Object, error) {
	return nil, errors.New("StateObject not implemented")
}

// StateObjects gets the result set for the specified key and state type.
func (coord *ServiceBase) StateObjects(t coordinator.StateType) (store.ResultSet, error) {
	return nil, errors.New("StateObjects not implemented")
}

// PostMessage posts the specified message to the coordinator.
func (coord *ServiceBase) PostMessage(msg coordinator.Message) error {
	return errors.New("PostMessage not implemented")
}

// SetNodeState posts the specified node status to the coordinator.
func (coord *ServiceBase) SetNodeState(node cluster.Node) error {
	return errors.New("SetNodeState not implemented")
}

// ClusterState gets the current cluster state.
func (coord *ServiceBase) ClusterState(name string) (cluster.Cluster, error) {
	return nil, errors.New("ClusterState not implemented")
}
