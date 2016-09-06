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
	"github.com/fanux/pbrain/plugins/pipeline"
	"github.com/spf13/cobra"
)

var (
	ManagerHost string
	ManagerPort string
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "A pipeline plugin",
	Long:  `With a simple config file, you can define when、what、how many your app run`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("pipeline called")
		RunPlugin(&pipeline.Pipeline{common.GetBasePlugin(ManagerHost, ManagerPort, pipeline.PLUGIN_NAME), nil})
	},
}

func init() {
	RootCmd.AddCommand(pipelineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipelineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pipelineCmd.Flags().StringVarP(&ManagerHost, "manager-host", "H", "localhost", "the ip address of plugin manager server")
	pipelineCmd.Flags().StringVarP(&ManagerPort, "manager-port", "p", ":8081", "the port of plugin manager server")
}
