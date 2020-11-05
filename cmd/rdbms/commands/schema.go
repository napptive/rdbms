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
	"time"

	"github.com/napptive/rdbms/internal/app/schema"

	"github.com/spf13/cobra"
)

var schemaCmdLongHelp = "Commands related with the database schema"
var schemaCmdShortHelp = "Commands related with the database schema"
var schemaCmdExample = `$ rdbms schema`
var schemaCmdUse = "schema"

var schemaCmd = &cobra.Command{
	Use:     schemaCmdUse,
	Long:    schemaCmdLongHelp,
	Example: schemaCmdExample,
	Short:   schemaCmdShortHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var loadSchemaCmdLongHelp = "Commands related with the database schema"
var loadSchemaCmdShortHelp = "Commands related with the database schema"
var loadSchemaCmdExample = `$ rdbms schema`
var loadSchemaCmdUse = "schema"

var loadSchemaVarFilepath string
var loadSchemaVarDuration time.Duration

var loadSchemaCmd = &cobra.Command{
	Use:     loadSchemaCmdUse,
	Long:    loadSchemaCmdLongHelp,
	Example: loadSchemaCmdExample,
	Short:   loadSchemaCmdShortHelp,

	RunE: func(cmd *cobra.Command, args []string) error {
		return schema.Load(loadSchemaVarFilepath, loadSchemaVarDuration, cfg)
	},
}

func init() {
	schemaCmd.PersistentFlags().StringVar(&loadSchemaVarFilepath, "scriptLoadPath", "", "Path where the load sql script is located.")
	schemaCmd.PersistentFlags().DurationVar(&loadSchemaVarDuration, "defaultStepTimeout", 30*time.Second, "Default time for each SQL script step")
	rootCmd.AddCommand(schemaCmd)
}