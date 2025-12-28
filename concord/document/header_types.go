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

const (
	V1 = KeyVersion(1)
)

const (
	CBOR = DocumentType(1)
)

const (
	StateHeaderObject   = HeaderType('S')
	MessageHeaderObject = HeaderType('M')
)

const (
	PrimaryIndex   = IndexType(1)
	SecondaryIndex = IndexType(2)
)

var (
	StateObjectKeyHeader   = [2]byte{byte(StateHeaderObject), byte(byte(CBOR) | HeaderByteFromVersion(V1))}
	MessageObjectKeyHeader = [2]byte{byte(MessageHeaderObject), byte(byte(CBOR) | HeaderByteFromVersion(V1))}
)

func HeaderByteFromVersion(v KeyVersion) byte {
	return (byte(v<<4) & 0x70)
}

func VertionFromHeaderByte(b byte) KeyVersion {
	return KeyVersion((b >> 4) & 0x07)
}

func TypeFromHeaderByte(b byte) byte {
	return (b & 0x07)
}

// HeaderTypes returns all header types.
func HeaderTypes() []HeaderType {
	return []HeaderType{
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
