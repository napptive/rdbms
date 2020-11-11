/**
 * Copyright 2020 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/napptive/rdbms/internal/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var cfg config.Config

var debugLevel bool
var consoleLogging bool

var rootCmdLongHelp = "This command contail useful operations to manage a Postgress database."
var rootCmdShortHelp = "RDBMS command"
var rootCmdExample = `
  $ rdbms help
  $ rdbms ping -c "host=localhost user=postgres password=Pass2020! port=5432"
  $ rdbms schema load --scriptLoadPath test/data/ValidSQLScript.yaml
  $ rdbms schema load --scriptLoadPath test/data/ValidSQLScript.yaml --selectedStep creation-step --selectedStep drop-step 
`
var rootCmdUse = "rdbms"

var rootCmd = &cobra.Command{
	Use:     rootCmdUse,
	Example: rootCmdExample,
	Short:   rootCmdShortHelp,
	Long:    rootCmdLongHelp,
	Version: "NaN",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debugLevel, "debug", "d", false, "Set debug level")
	rootCmd.PersistentFlags().BoolVarP(&consoleLogging, "consoleLogging", "l", false, "Pretty print logging")

	rootCmd.PersistentFlags().StringVarP(&cfg.ConnString, "connectionString", "c",
		"host=localhost user=postgres password=Pass2020! port=5432",
		"Database connection string")

	rootCmd.PersistentFlags().IntVarP(&cfg.PingRetries, "pingRetries", "r", 3,
		"Number of retries to ping to the database.")
	rootCmd.PersistentFlags().DurationVarP(&cfg.PingWaitingPeriod, "pingWaitingPeriod", "w",
		5*time.Second, "Waiting time between each ping command")

	rootCmd.PersistentFlags().BoolVarP(&cfg.SkipPing, "skipPing", "k",
		false, "If true, the command skip the ping step.")
}

// Execute the user command
func Execute(version string, commit string) {
	versionTemplate := fmt.Sprintf("%s [%s] ", version, commit)
	rootCmd.SetVersionTemplate(versionTemplate)
	cfg.Version = version
	cfg.Commit = commit

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	setupLogging()
}

// setupLogging sets the debug level and console logging if required.
func setupLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if consoleLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
}
