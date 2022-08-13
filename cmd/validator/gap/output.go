// Copyright © 2021 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validatorgap

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type dataOut struct {
	debug   bool
	quiet   bool
	verbose bool
}

func output(ctx context.Context, data *dataOut) (string, error) {
	if data == nil {
		return "", errors.New("no data")
	}

	if data.quiet {
		return "", nil
	}

	builder := strings.Builder{}

	now := time.Now()
	builder.WriteString("Current time: ")
	builder.WriteString(now.Format("15:04:05\n"))

	// (Validator: 0xdeadbeef)
	// (Upcoming attestation slot: 09:51:47 in this/next epoch)
	// [Upcoming proposer slot: 10:41:28 in this epoch]
	//
	// (Validator: 0xcafebabe)
	// (Upcoming attestation slot: 09:51:47 in this/next epoch)
	// [Upcoming proposer slot: 10:41:28 in this epoch]
	//
	// Longest gap: 00:10:15 - 00:15:15 (600s) in epoch 12345

	return builder.String(), nil
}
