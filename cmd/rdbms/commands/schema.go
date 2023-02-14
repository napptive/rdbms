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

	"github.com/napptive/rdbms/v2/internal/app/rdbms"

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

var loadSchemaCmdLongHelp = "Load execute a set of queries in the target database"
var loadSchemaCmdShortHelp = "Load command"
var loadSchemaCmdExample = `
  $ rdbms schema load --scriptLoadPath test/data/ValidSQLScript.yaml
  $ rdbms schema load --scriptLoadPath test/directory-with-scripts/ --fileExtension .sql
  $ rdbms schema load --scriptLoadPath test/data/ValidSQLScript.yaml --selectedStep creation-step --selectedStep drop-step 
`
var loadSchemaCmdUse = "load"

var loadSchemaVarFilepath string
var loadSchemaFileExtension string
var loadSchemaVarDuration time.Duration
var loadSchemaVarSteps []string

var loadSchemaCmd = &cobra.Command{
	Use:     loadSchemaCmdUse,
	Long:    loadSchemaCmdLongHelp,
	Example: loadSchemaCmdExample,
	Short:   loadSchemaCmdShortHelp,

	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := rdbms.Load(loadSchemaVarFilepath, loadSchemaFileExtension, loadSchemaVarDuration, loadSchemaVarSteps, cfg)
		if err != nil {
			return err
		}
		for _, r := range result {
			r.Print()
		}
		return nil
	},
}

func init() {
	loadSchemaCmd.PersistentFlags().StringVar(&loadSchemaVarFilepath, "scriptLoadPath", "", "Path where the load sql script is located. This could be a single file or a directory")
	err := loadSchemaCmd.MarkPersistentFlagRequired("scriptLoadPath")
	if err != nil {
		panic(err)
	}
	loadSchemaCmd.Flags().StringVar(&loadSchemaFileExtension, "fileExtension", ".yaml", "File extension that is applied to the listing of files when scriptLoadPath points to a directory.")
	loadSchemaCmd.PersistentFlags().DurationVar(&loadSchemaVarDuration, "defaultStepTimeout", 30*time.Second, "Default time for each SQL script step")
	loadSchemaCmd.PersistentFlags().StringArrayVar(&loadSchemaVarSteps, "selectedStep", []string{}, "Select the steps that you want to execute, if empty you execute all the steps.")
	schemaCmd.AddCommand(loadSchemaCmd)

	rootCmd.AddCommand(schemaCmd)
}
