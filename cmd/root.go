// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thales-e-security/contribstats/pkg/config"
	"github.com/thales-e-security/contribstats/pkg/server"
)

var cfgFile string
var debug bool
var s *server.StatServer

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "contribstats",
	Short: "Collect GitHub Stats for an Organization",
	Long:  `Collect GitHub Stats for an Organization`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		s = server.NewStatServer()
		if err := s.Start(); err != nil {
			logrus.Panic(err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(func() {
	//	cfgFile = config.InitConfig(cfgFile)
	//})
	cfgFile = config.InitConfig(cfgFile)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.contribstats.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug Logging")
}
