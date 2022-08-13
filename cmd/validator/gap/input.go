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
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type dataIn struct {
	// System.
	timeout time.Duration
	quiet   bool
	verbose bool
	debug   bool
	// Ethereum 2 connection.
	eth2Client    string
	allowInsecure bool
	// Operation.
	accounts []string
	pubkeys  []string
	indices  []string
}

func input(ctx context.Context) (*dataIn, error) {
	data := &dataIn{}

	if viper.GetDuration("timeout") == 0 {
		return nil, errors.New("timeout is required")
	}
	data.timeout = viper.GetDuration("timeout")
	data.quiet = viper.GetBool("quiet")
	data.verbose = viper.GetBool("verbose")
	data.debug = viper.GetBool("debug")

	// Ethereum 2 connection.
	data.eth2Client = viper.GetString("connection")
	data.allowInsecure = viper.GetBool("allow-insecure-connections")

	// Validators.
	data.indices = viper.GetStringSlice("accounts")
	data.indices = viper.GetStringSlice("pubkeys")
	data.indices = viper.GetStringSlice("indices")

	if len(data.accounts) == 0 && len(data.pubkeys) == 0 && len(data.indices) == 0 {
		return nil, errors.New("a list of accounts, pubkeys or indices is required")
	}

	return data, nil
}
