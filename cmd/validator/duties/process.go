// Copyright © 2019 - 2022 Weald Technology Trading
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

package validatorduties

import (
	"context"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
	"github.com/wealdtech/ethdo/util"
)

func process(ctx context.Context, data *dataIn) (*dataOut, error) {
	if data == nil {
		return nil, errors.New("no data")
	}

	// Ethereum 2 client.
	eth2Client, err := util.ConnectToBeaconNode(ctx, data.eth2Client, data.timeout, data.allowInsecure)
	if err != nil {
		return nil, err
	}

	results := &dataOut{
		debug:   data.debug,
		quiet:   data.quiet,
		verbose: data.verbose,
	}

	validatorIndex, err := util.ValidatorIndex(ctx, eth2Client, data.account, data.pubKey, data.index)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validator index")
	}

	// Fetch duties for this and next epoch.
	thisEpoch, err := util.CurrentEpoch(ctx, eth2Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate current epoch")
	}
	thisEpochAttesterDuty, err := util.AttesterDuty(ctx, eth2Client, validatorIndex, thisEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain this epoch attester duty for validator")
	}
	results.thisEpochAttesterDuty = thisEpochAttesterDuty

	thisEpochProposerDuties, err := util.ProposerDuties(ctx, eth2Client, validatorIndex, thisEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain this epoch proposer duties for validator")
	}
	results.thisEpochProposerDuties = thisEpochProposerDuties

	nextEpoch := thisEpoch + 1
	nextEpochAttesterDuty, err := util.AttesterDuty(ctx, eth2Client, validatorIndex, nextEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain next epoch attester duty for validator")
	}
	results.nextEpochAttesterDuty = nextEpochAttesterDuty

	genesis, err := eth2Client.(eth2client.GenesisProvider).Genesis(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain genesis data")
	}
	results.genesisTime = genesis.GenesisTime

	config, err := eth2Client.(eth2client.SpecProvider).Spec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain beacon chain configuration")
	}
	results.slotsPerEpoch = config["SLOTS_PER_EPOCH"].(uint64)
	results.slotDuration = config["SECONDS_PER_SLOT"].(time.Duration)

	return results, nil
}
