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
	"errors"

	"github.com/cybergarage/go-concord/concord"
)

func truncateCoordinatorStore(coord concord.Service) error {
	txn, err := coord.Transact()
	if err != nil {
		return err
	}
	err = txn.Truncate()
	if err != nil {
		return errors.Join(err, txn.Cancel())
	}
	return txn.Commit()
}
