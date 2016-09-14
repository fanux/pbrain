// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
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

	"github.com/fanux/pbrain/common"
	"github.com/fanux/pbrain/plugins/decider"
	"github.com/spf13/cobra"
)

// deciderCmd represents the decider command
var deciderCmd = &cobra.Command{
	Use:   "decider",
	Short: "scale apps by app metrical",
	Long:  `app metrical is loadbalance info, or cpu memery use info`,
	Run: func(cmd *cobra.Command, args []string) {
		//  Work your own magic here
		fmt.Println("decider called")
		basePlugin := common.GetBasePlugin(ManagerHost, ManagerPort, decider.PLUGIN_NAME)
		RunPlugin(&decider.Decider{basePlugin, nil})
	},
}

func init() {
	RootCmd.AddCommand(deciderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deciderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deciderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
