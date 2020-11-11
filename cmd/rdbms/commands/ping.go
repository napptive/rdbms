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
	"github.com/napptive/rdbms/internal/app/rdbms"
	"github.com/spf13/cobra"
)

var pingCmdLongHelp = "Check if the database is alived."
var pingCmdShortHelp = "Ping databases"
var pingCmdExample = `$ rdbms ping`
var pingCmdUse = "ping"

var pingCmd = &cobra.Command{
	Use:     pingCmdUse,
	Long:    pingCmdLongHelp,
	Example: pingCmdExample,
	Short:   pingCmdShortHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rdbms.Ping(cfg)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
