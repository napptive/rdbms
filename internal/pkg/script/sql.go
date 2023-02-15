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

package script

import (
	"fmt"
	"os"
	"time"

	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// SQL represents the import file to execute batch command in the database.
type SQL struct {
	Steps []SQLStep `yaml:"steps,flow"`
}

// SQLStep represent a set of setence that must executed in the same transaction.
type SQLStep struct {
	Name    string   `yaml:"name"`
	Timeout string   `yaml:"timeout,omitempty"`
	Queries []string `yaml:"queries,flow"`
}

// TimeoutDuration parses the time duration of the step.
func (step *SQLStep) TimeoutDuration(defaultTimeout time.Duration) (time.Duration, error) {
	if step.Timeout == "" {
		return defaultTimeout, nil
	}
	return time.ParseDuration(step.Timeout)
}

// SQLFileParse transform a file in a script SQL struct.
func SQLFileParse(filepath string) (*SQL, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return SQLParse(data)
}

func SQLParse(data []byte) (*SQL, error) {
	var result SQL
	m := yaml.MapSlice{}
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	if len(m) <= 0 {
		return nil, nerrors.NewInternalError("Invalid data in yaml file")
	}
	// check entry key
	if m[0].Key == "steps" {

		steps := m[0].Value.([]interface{})
		for _, s := range steps {
			r := s.(yaml.MapSlice)
			var sql SQLStep
			for _, entry := range r {
				switch entry.Key {
				case "name":
					sql.Name = fmt.Sprintf("%v", entry.Value)
				case "timeout":
					sql.Timeout = fmt.Sprintf("%v", entry.Value)
				case "queries":
					query := entry.Value.([]interface{})
					for _, q := range query {
						sql.Queries = append(sql.Queries, fmt.Sprintf("%v", q))
					}
				default:
					return nil, nerrors.NewInternalError("unexpected key found in yaml [%s]", entry.Key)
				}
			}
			result.Steps = append(result.Steps, sql)
			log.Debug().Interface("sql", sql).Msg("SQL!")
		}
	}
	return &result, nil
}
