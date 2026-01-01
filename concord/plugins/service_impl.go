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
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/cybergarage/go-cbor/cbor"
	"github.com/cybergarage/go-concord/concord"
	"github.com/cybergarage/go-concord/concord/cluster"
	"github.com/cybergarage/go-concord/concord/coordinator"
	"github.com/cybergarage/go-concord/concord/document"
	"github.com/cybergarage/go-concord/concord/store"
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-serix/serix/plugins/document/key/composite"
)

const (
	DefaultStoreScanInterval = time.Second
)

// KeyCoder represents a key coder for document keys.
type KeyCoder = store.KeyCoder

// Store represents a coordinator store.
type Store = coordinator.Store

type serviceImpl struct {
	name string
	Store
	cluster.Node
	*time.Ticker
	observers []coordinator.Observer
	*coordinator.MessageQueue
	ctx       context.Context
	ctxCancel context.CancelFunc
}

// NewDefaultKeyCoder returns a new default key coder for the coordinator service.
func NewDefaultKeyCoder() store.KeyCoder {
	return composite.NewCoder()
}

// ServiceOptions represents options for the coordinator service.
type ServiceOptions func(service *serviceImpl)

// WithServiceName returns a service option with the specified service name.
func WithServiceName(name string) ServiceOptions {
	return func(service *serviceImpl) {
		service.name = name
	}
}

// WithServiceStore returns a service option with the specified store.
func WithServiceStore(store store.Store) ServiceOptions {
	return func(service *serviceImpl) {
		service.Store = coordinator.NewStoreWith(store)
	}
}

// NewService returns a new coordinator service instance.
func NewService(opts ...ServiceOptions) (Service, error) {
	ctx, cancel := context.WithCancel(context.Background())
	coord := &serviceImpl{
		name:         "",
		Store:        nil,
		Node:         cluster.NewNode(),
		Ticker:       time.NewTicker(time.Second),
		observers:    make([]concord.Observer, 0),
		MessageQueue: coordinator.NewMessageQueue(),
		ctx:          ctx,
		ctxCancel:    cancel,
	}
	for _, opt := range opts {
		opt(coord)
	}
	if coord.Store == nil {
		return nil, errors.New("store is not set")
	}
	return coord, nil
}

// ServiceName returns the name of the coordinator service.
func (coord *serviceImpl) ServiceName() string {
	return coord.name
}

// AddObserver adds the specified observer.
func (coord *serviceImpl) AddObserver(newObserver coordinator.Observer) error {
	for _, observer := range coord.observers {
		if observer == newObserver {
			return nil
		}
	}
	coord.observers = append(coord.observers, newObserver)
	return nil
}

// SetStateObject sets the state object for the specified key.
func (coord *serviceImpl) SetStateObject(t coordinator.StateType, obj coordinator.Object) error {
	txn, err := coord.Transact()
	if err != nil {
		return err
	}
	steteKey := coordinator.NewStateKeyWith(t, obj.Key()...)
	err = txn.Set(document.NewObjectWith(steteKey, obj.Bytes()))
	if err != nil {
		return errors.Join(err, txn.Cancel())
	}
	return txn.Commit()
}

// StateObject gets the object for the specified key and state type.
func (coord *serviceImpl) StateObject(t coordinator.StateType, key coordinator.Key) (coordinator.Object, error) {
	txn, err := coord.Transact()
	if err != nil {
		return nil, err
	}
	steteKey := coordinator.NewStateKeyWith(t, key...)
	obj, err := txn.Get(steteKey)
	if err != nil {
		return nil, errors.Join(err, txn.Cancel())
	}
	err = txn.Commit()
	return obj, err
}

// StateObjects gets the result set for the specified key and state type.
func (coord *serviceImpl) StateObjects(t coordinator.StateType) (coordinator.ResultSet, error) {
	txn, err := coord.Transact()
	if err != nil {
		return nil, err
	}
	rs, err := txn.Scan(coordinator.NewScanStateKeyWith(t))
	if err != nil {
		return nil, errors.Join(err, txn.Cancel())
	}
	err = txn.Commit()
	return rs, err
}

// nofityMessage posts the specified message to the observers.
func (coord *serviceImpl) nofityMessage(msg coordinator.Message) {
	for _, observer := range coord.observers {
		observer.OnMessageReceived(msg)
	}
}

func (coord *serviceImpl) getLatestMessages(txn coordinator.Transaction) (coordinator.ResultSet, error) {
	key := coordinator.NewMessageScanKey()
	rs, err := txn.Scan(
		key,
		store.OrderDesc)
	return rs, err
}

func (coord *serviceImpl) notifyUpdateMessages(txn coordinator.Transaction) error {
	rs, err := coord.getLatestMessages(txn)
	if err != nil {
		return err
	}
	defer rs.Close()

	localClock := coord.Clock()

	msgs := []coordinator.Message{}
	for rs.Next() {
		msgObj := coordinator.NewMessageObject()
		r, err := rs.Result()
		if err != nil {
			return err
		}
		obj, err := r.Object()
		if err != nil {
			return err
		}
		err = obj.Unmarshal(msgObj)
		if err != nil {
			return err
		}

		// Skip the message if the message clock is older than the local clock
		if 0 < cluster.CompareClocks(localClock, msgObj.MsgClock) {
			break
		}

		// Skip the self messages
		if msgObj.FromHost == coord.Host() {
			continue
		}

		msg := coordinator.NewMessageFrom(msgObj)
		msgs = append([]coordinator.Message{msg}, msgs...)

		coord.SetReceivedClock(msgObj.MsgClock)
	}
	if err := rs.Err(); err != nil {
		return err
	}

	for _, msg := range msgs {
		log.Infof("RECV message: %s %s (%d)", msg.From().Host(), msg.Event().String(), msg.Clock())
		coord.nofityMessage(msg)
	}

	return nil
}

func (coord *serviceImpl) getLatestMessageClock(txn coordinator.Transaction) (cluster.Clock, error) {
	rs, err := coord.getLatestMessages(txn)
	if err != nil {
		return 0, err
	}
	defer rs.Close()

	if !rs.Next() {
		return 0, rs.Err()
	}

	msgObj := coordinator.NewMessageObject()
	r, err := rs.Result()
	if err != nil {
		return 0, err
	}
	obj, err := r.Object()
	if err != nil {
		return 0, err
	}
	err = obj.Unmarshal(msgObj)
	if err != nil {
		return 0, err
	}

	return msgObj.MsgClock, nil
}

// PostMessage posts the specified message to the coordinator.
func (coord *serviceImpl) PostMessage(msg coordinator.Message) error {
	coord.Lock()
	defer coord.Unlock()

	coord.EnqueueMessage(msg)

	return nil
}

// postMessage posts the specified message to the coordinator.
func (coord *serviceImpl) postMessage(txn coordinator.Transaction, msg coordinator.Message) error {
	localClock := coord.IncrementClock()

	obj, err := coordinator.NewMessageObjectWith(msg, coord, localClock)
	if err != nil {
		return errors.Join(err, txn.Cancel())
	}

	objBytes, err := cbor.Marshal(obj)
	if err != nil {
		return err
	}

	log.Infof("SEND message: %s %s (%d)", obj.FromHost, msg.Event().String(), obj.MsgClock)

	key := coordinator.NewMessageKeyWith(msg, localClock)
	err = txn.Set(document.NewObjectWith(key, objBytes))
	if err != nil {
		return err
	}

	return nil
}

func (coord *serviceImpl) postNodeState(txn coordinator.Transaction, node cluster.Node) error {
	key := coordinator.NewNodeKeyWith(node)
	obj := coordinator.NewNodeObjectWith(node)
	objBytes, err := cbor.Marshal(obj)
	if err != nil {
		return err
	}

	err = txn.Set(document.NewObjectWith(key, objBytes))
	if err != nil {
		return err
	}

	return nil
}

// SetNodeState posts the specified node state to the coordinator.
func (coord *serviceImpl) SetNodeState(node cluster.Node) error {
	coord.Lock()
	defer coord.Unlock()

	txn, err := coord.Transact()
	if err != nil {
		return err
	}

	err = coord.postNodeState(txn, node)
	if err != nil {
		return errors.Join(err, txn.Cancel())
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

// ClusterState gets the current cluster state.
func (coord *serviceImpl) ClusterState(name string) (cluster.Cluster, error) {
	coord.Lock()
	defer coord.Unlock()

	txn, err := coord.Transact()
	if err != nil {
		return nil, err
	}

	rs, err := txn.Scan(coordinator.NewClusterScanKeyWith(name))
	if err != nil {
		return nil, errors.Join(err, txn.Cancel())
	}
	defer rs.Close()

	nodes := []cluster.Node{}
	for rs.Next() {
		nodeObj := coordinator.NewNodeObject()
		r, err := rs.Result()
		if err != nil {
			return nil, errors.Join(err, txn.Cancel())
		}
		obj, err := r.Object()
		if err != nil {
			return nil, errors.Join(err, txn.Cancel())
		}
		err = obj.Unmarshal(nodeObj)
		if err != nil {
			return nil, errors.Join(err, txn.Cancel())
		}

		node, err := coordinator.NewNodeWith(nodeObj)
		if err != nil {
			return nil, errors.Join(err, txn.Cancel())
		}
		nodes = append(nodes, node)
	}
	if err := rs.Err(); err != nil {
		return nil, errors.Join(err, txn.Cancel())
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}

	return cluster.NewClusterWith(name, nodes), nil
}

// Start starts this etcd coordinator.
func (coord *serviceImpl) Start() error { // nolint:gocognit
	if err := coord.Store.Start(); err != nil {
		return err
	}

	txn, err := coord.Transact()
	if err != nil {
		return err
	}

	// Set latest message clock to the local clock

	clock, err := coord.getLatestMessageClock(txn)
	if err != nil {
		return errors.Join(err, txn.Cancel())
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	coord.SetClock(clock)
	coord.IncrementClock()

	// Start coordinator worker

	go func() {
		logError := func(err error) {
			log.Warnf("coordinator worker: %s", err)
		}

		pushPostedMessages := func(postedMsgs []coordinator.Message) {
			for n := len(postedMsgs) - 1; 0 <= n; n-- {
				coord.PushMessage(postedMsgs[n])
			}
		}

		for {
			jitter := time.Duration(rand.Intn(int(DefaultStoreScanInterval/time.Millisecond/2))) * time.Millisecond //nolint:gosec
			select {
			case <-time.After(DefaultStoreScanInterval + jitter):
				var err error
				coord.Lock()

				startClock := coord.Clock()

				// Start transaction

				txn, err := coord.Transact()
				if err != nil {
					logError(err)
					coord.Unlock()
					continue
				}

				// Receive update messages and update local clock

				err = coord.notifyUpdateMessages(txn)
				if err != nil {
					logError(errors.Join(err, txn.Cancel()))
					coord.Unlock()
					continue
				}

				// Post message if there is no message in the queue

				postedMsgs := []coordinator.Message{}
				msg, err := coord.PopMessage()
				for msg != nil && err == nil {
					err = coord.postMessage(txn, msg)
					if err != nil {
						coord.PushMessage(msg)
						pushPostedMessages(postedMsgs)
						break
					}
					postedMsgs = append(postedMsgs, msg)
					msg, err = coord.PopMessage()
				}

				if err != nil && !errors.Is(err, coordinator.ErrNoMessage) {
					logError(errors.Join(err, txn.Cancel()))
					coord.Unlock()
					continue
				}

				// Update node status

				if 0 < cluster.CompareClocks(coord.Clock(), startClock) {
					err := coord.postNodeState(txn, coord)
					if err != nil {
						logError(errors.Join(err, txn.Cancel()))
					}
				}

				// Commit transaction

				err = txn.Commit()
				if err != nil {
					pushPostedMessages(postedMsgs)
					logError(err)
				}

				coord.Unlock()
			case <-coord.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Stop stops this etcd coordinator.
func (coord *serviceImpl) Stop() error {
	coord.ctxCancel()
	<-coord.ctx.Done()

	if err := coord.Store.Stop(); err != nil {
		return err
	}

	return nil
}
