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

package document

import (
	"github.com/cybergarage/go-serix/serix/document/kv"
)

const (
	V1 = Version(1)
)

const (
	CBOR = Format(1)
)

const (
	StateHeaderObject   = Category('S')
	MessageHeaderObject = Category('M')
)

// KeyHeader represents a header for all keys.
type KeyHeader = kv.KeyHeader

// Category represents a category.
type Category = kv.Category

// Version represents a version.
type Version = kv.Version

// Format represents a format.
type Format = kv.Format

// IndexType represents an index type.
type IndexType byte

// NewKeyHeaderFrom creates a new key header from the specified bytes.
func NewKeyHeaderFrom(b []byte) KeyHeader {
	var header KeyHeader
	copy(header[:], b)
	return header
}

var (
	StateObjectKeyHeader   = [2]byte{byte(StateHeaderObject), byte(byte(CBOR) | V1.HeaderByte())}
	MessageObjectKeyHeader = [2]byte{byte(MessageHeaderObject), byte(byte(CBOR) | V1.HeaderByte())}
)

// Categorys returns all header types.
func Categorys() []Category {
	return []Category{
		StateHeaderObject,
		MessageHeaderObject,
	}
}

// HeaderPrefixes returns all header prefixes.
func HeaderPrefixes() [][]byte {
	return [][]byte{
		StateObjectKeyHeader[:],
		MessageObjectKeyHeader[:],
	}
}

// HeaderPrefixKeys returns all header prefix keys.
func HeaderPrefixKeys() []Key {
	keys := []Key{}
	for _, prefix := range HeaderPrefixes() {
		key := NewKeyWith(prefix)
		keys = append(keys, key)
	}
	return keys
}
