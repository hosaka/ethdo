// Copyright © 2022 Weald Technology Trading
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

package util

import (
	"context"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

func AttesterDuty(ctx context.Context, eth2Client eth2client.Service, validatorIndex spec.ValidatorIndex, epoch spec.Epoch) (*api.AttesterDuty, error) {
	// Find the attesting slot for the given epoch.
	duties, err := eth2Client.(eth2client.AttesterDutiesProvider).AttesterDuties(ctx, epoch, []spec.ValidatorIndex{validatorIndex})
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain attester duties")
	}

	if len(duties) == 0 {
		return nil, errors.New("validator does not have duty for that epoch")
	}

	return duties[0], nil
}

func ProposerDuties(ctx context.Context, eth2Client eth2client.Service, validatorIndex spec.ValidatorIndex, epoch spec.Epoch) ([]*api.ProposerDuty, error) {
	// Fetch the proposer duties for this epoch.
	proposerDuties, err := eth2Client.(eth2client.ProposerDutiesProvider).ProposerDuties(ctx, epoch, []spec.ValidatorIndex{validatorIndex})
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain proposer duties")
	}

	return proposerDuties, nil
}

func CurrentEpoch(ctx context.Context, eth2Client eth2client.Service) (spec.Epoch, error) {
	config, err := eth2Client.(eth2client.SpecProvider).Spec(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain beacon chain configuration")
	}
	slotsPerEpoch := config["SLOTS_PER_EPOCH"].(uint64)
	slotDuration := config["SECONDS_PER_SLOT"].(time.Duration)
	genesis, err := eth2Client.(eth2client.GenesisProvider).Genesis(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain genesis data")
	}

	if genesis.GenesisTime.After(time.Now()) {
		return spec.Epoch(0), nil
	}
	return spec.Epoch(uint64(time.Since(genesis.GenesisTime).Seconds()) / (uint64(slotDuration.Seconds()) * slotsPerEpoch)), nil
}
