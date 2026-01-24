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
	"testing"
)

func TestKeyHeader(t *testing.T) {
	type expected struct {
		category Category
		ver      Version
		fmt      Format
	}
	testKeyHeaders := []struct {
		header   KeyHeader
		expected expected
	}{
		{
			header: StateObjectKeyHeader,
			expected: expected{
				category: StateHeaderObject,
				ver:      V1,
				fmt:      CBOR,
			},
		},
		{
			header: MessageObjectKeyHeader,
			expected: expected{
				category: MessageHeaderObject,
				ver:      V1,
				fmt:      CBOR,
			},
		},
	}
	for _, key := range testKeyHeaders {
		if key.header.Category() != key.expected.category {
			t.Errorf("%v != %v", key.header.Category(), key.expected.category)
		}
		if key.header.Version() != key.expected.ver {
			t.Errorf("%v != %v", key.header.Version(), key.expected.ver)
		}
		if key.expected.fmt != Format(0) {
			if key.header.Format() != key.expected.fmt {
				t.Errorf("%v != %v", key.header.Format(), key.expected.fmt)
			}
		}
	}
}
