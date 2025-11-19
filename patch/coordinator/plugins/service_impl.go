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
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/cybergarage/go-cbor/cbor"
	"github.com/cybergarage/go-concord/concord/cluster"
	"github.com/cybergarage/go-concord/concord/core"
	"github.com/cybergarage/go-logger/log"
)

const (
	DefaultStoreScanInterval = time.Second
)

type serviceImpl struct {
	Coordinator
	observers []Observer
	*core.MessageQueue
	ctx       context.Context
	ctxCancel context.CancelFunc
}

// NewServiceWith returns a new coordinator service with the specified core coordinator service.
func NewServiceWith(coord Coordinator) Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &serviceImpl{
		Coordinator:  coord,
		observers:    make([]Observer, 0),
		MessageQueue: core.NewMessageQueue(),
		ctx:          ctx,
		ctxCancel:    cancel,
	}
}

// AddObserver adds the specified observer.
func (coord *serviceImpl) AddObserver(newObserver Observer) error {
	for _, observer := range coord.observers {
		if observer == newObserver {
			return nil
		}
	}
	coord.observers = append(coord.observers, newObserver)
	return nil
}

// SetStateObject sets the state object for the specified key.
func (coord *serviceImpl) SetStateObject(t StateType, obj Object) error {
	txn, err := coord.Transact()
	if err != nil {
		return err
	}
	stateKey := core.NewStateKeyWith(t, obj.Key()...)
	err = txn.Set(core.NewObjectWith(stateKey, obj.Bytes()))
	if err != nil {
		return errors.Join(err, txn.Cancel())
	}
	return txn.Commit()
}

// GetStateObject gets the state object for the specified key and state type.
func (coord *serviceImpl) GetStateObject(t StateType, key Key) (Object, error) {
	txn, err := coord.Transact()
	if err != nil {
		return nil, err
	}
	stateKey := core.NewStateKeyWith(t, key...)
	obj, err := txn.Get(stateKey)
	if err != nil {
		return nil, errors.Join(err, txn.Cancel())
	}
	err = txn.Commit()
	return obj, err
}

// GetStateObjects gets the result set for the specified key and state type.
func (coord *serviceImpl) GetStateObjects(t StateType) (ResultSet, error) {
	txn, err := coord.Transact()
	if err != nil {
		return nil, err
	}
	rs, err := txn.GetRange(core.NewScanStateKeyWith(t))
	if err != nil {
		return nil, errors.Join(err, txn.Cancel())
	}
	err = txn.Commit()
	return rs, err
}

// nofityMessage posts the specified message to the observers.
func (coord *serviceImpl) nofityMessage(msg Message) {
	for _, observer := range coord.observers {
		observer.OnMessageReceived(msg)
	}
}

func (coord *serviceImpl) getLatestMessages(txn Transaction) (ResultSet, error) {
	key := core.NewMessageScanKey()
	rs, err := txn.GetRange(
		key,
		core.NewOrderOptionWith(core.OrderDesc))
	return rs, err
}

func (coord *serviceImpl) notifyUpdateMessages(txn Transaction) error {
	rs, err := coord.getLatestMessages(txn)
	if err != nil {
		return err
	}

	localClock := coord.Clock()

	msgs := []Message{}
	for rs.Next() {
		msgObj := core.NewMessageObject()
		obj := rs.Object()
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

		msg := core.NewMessageFrom(msgObj)
		msgs = append([]Message{msg}, msgs...)

		coord.SetReceivedClock(msgObj.MsgClock)
	}

	for _, msg := range msgs {
		log.Infof("RECV message: %s %s (%d)", msg.From().Host(), msg.Event().String(), msg.Clock())
		coord.nofityMessage(msg)
	}

	return nil
}

func (coord *serviceImpl) getLatestMessageClock(txn Transaction) (cluster.Clock, error) {
	rs, err := coord.getLatestMessages(txn)
	if err != nil {
		return 0, err
	}

	if !rs.Next() {
		return 0, nil
	}

	msgObj := core.NewMessageObject()
	obj := rs.Object()
	err = obj.Unmarshal(msgObj)
	if err != nil {
		return 0, err
	}

	return msgObj.MsgClock, nil
}

// PostMessage posts the specified message to the
func (coord *serviceImpl) PostMessage(msg Message) error {
	coord.Lock()
	defer coord.Unlock()

	coord.EnqueueMessage(msg)

	return nil
}

// postMessage posts the specified message to the
func (coord *serviceImpl) postMessage(txn Transaction, msg Message) error {
	localClock := coord.IncrementClock()

	obj, err := core.NewMessageObjectWith(msg, coord, localClock)
	if err != nil {
		return errors.Join(err, txn.Cancel())
	}

	objBytes, err := cbor.Marshal(obj)
	if err != nil {
		return err
	}

	log.Infof("SEND message: %s %s (%d)", obj.FromHost, msg.Event().String(), obj.MsgClock)

	key := core.NewMessageKeyWith(msg, localClock)
	err = txn.Set(core.NewObjectWith(key, objBytes))
	if err != nil {
		return err
	}

	return nil
}

func (coord *serviceImpl) postNodeState(txn Transaction, node cluster.Node) error {
	key := core.NewNodeKeyWith(node)
	obj := core.NewNodeObjectWith(node)
	objBytes, err := cbor.Marshal(obj)
	if err != nil {
		return err
	}

	err = txn.Set(core.NewObjectWith(key, objBytes))
	if err != nil {
		return err
	}

	return nil
}

// SetNodeState posts the specified node state to the
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

// GetClusterState gets the current cluster state.
func (coord *serviceImpl) GetClusterState(name string) (cluster.Cluster, error) {
	coord.Lock()
	defer coord.Unlock()

	txn, err := coord.Transact()
	if err != nil {
		return nil, err
	}

	rs, err := txn.GetRange(core.NewClusterScanKeyWith(name))
	if err != nil {
		return nil, errors.Join(err, txn.Cancel())
	}

	nodes := []cluster.Node{}
	for rs.Next() {
		nodeObj := core.NewNodeObject()
		obj := rs.Object()
		err = obj.Unmarshal(nodeObj)
		if err != nil {
			return nil, errors.Join(err, txn.Cancel())
		}

		node, err := core.NewNodeWith(nodeObj)
		if err != nil {
			return nil, errors.Join(err, txn.Cancel())
		}
		nodes = append(nodes, node)
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}

	return cluster.NewClusterWith(name, nodes), nil
}

// Start starts this etcd
func (coord *serviceImpl) Start() error { // nolint:gocognit
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

		pushPostedMessages := func(postedMsgs []Message) {
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

				postedMsgs := []Message{}
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

				if err != nil && !errors.Is(err, core.ErrNoMessage) {
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

// Stop stops this etcd
func (coord *serviceImpl) Stop() error {
	coord.ctxCancel()
	<-coord.ctx.Done()
	return nil
}
