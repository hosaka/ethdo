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

	validatorIndex, err := util.ValidatorIndex(ctx, eth2Client, "", "", data.indices[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validator index")
	}

	// todo: get attester gaps for each validator
	//  find and overlap in gaps to propose the longest gap
	//  check that proposer duties do not overlap with longest gap

	// Fetch duties for this and next epoch.
	thisEpoch, err := util.CurrentEpoch(ctx, eth2Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate current epoch")
	}
	thisEpochAttesterDuty, err := util.AttesterDuty(ctx, eth2Client, validatorIndex, thisEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain this epoch attester duty for validator")
	}

	thisEpochProposerDuties, err := util.ProposerDuties(ctx, eth2Client, validatorIndex, thisEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain this epoch proposer duties for validator")
	}

	nextEpoch := thisEpoch + 1
	nextEpochAttesterDuty, err := util.AttesterDuty(ctx, eth2Client, validatorIndex, nextEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain next epoch attester duty for validator")
	}

	println(thisEpochAttesterDuty.ValidatorIndex)
	if len(thisEpochProposerDuties) > 0 {
		println(thisEpochProposerDuties[0].ValidatorIndex)
	}
	println(nextEpochAttesterDuty.ValidatorIndex)

	return results, nil
}
