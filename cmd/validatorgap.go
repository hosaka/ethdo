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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	validatorgap "github.com/wealdtech/ethdo/cmd/validator/gap"
)

var validatorGapCmd = &cobra.Command{
	Use:   "gap",
	Short: "Calculate the longest time gap between validators duties",
	Long:  `Calculate the longest time gap between validator duties`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := validatorgap.Run(cmd)
		if err != nil {
			return err
		}
		if viper.GetBool("quiet") {
			return nil
		}
		fmt.Print(res)
		return nil
	},
}

func init() {
	validatorCmd.AddCommand(validatorGapCmd)
	validatorFlags(validatorGapCmd)
	validatorGapCmd.Flags().String("indices", "", "validator indices for duties gap, separated with spaces")
}

func validatorGapBindings() {
	if err := viper.BindPFlag("indices", validatorGapCmd.Flags().Lookup("indices")); err != nil {
		panic(err)
	}
}
