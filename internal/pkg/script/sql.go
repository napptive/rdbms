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
	"encoding/json"
	"io/ioutil"
	"time"
)



// SQL respresents the import file to execute batch command in the database.
type SQL struct {
	Steps []SQLStep `json:"steps"`
}

// SQLStep represent a set of setence that must executed in the same transaction.
type SQLStep struct {
	Name    string   `json:"name"`
	Timeout string   `json:"timeout"`
	Queries []string `json:"queries"`
}

//TimeoutDuration parses the time duration of the step.
func (step *SQLStep) TimeoutDuration(defaultTimeout time.Duration) (time.Duration, error) {
	if step.Timeout == "" {
		return defaultTimeout, nil
	}
	return time.ParseDuration(step.Timeout)
}

// SQLFileParse transform a file in a script SQL struct.
func SQLFileParse(filepath string) (*SQL, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return SQLParse(data)
}

// SQLParse transform a string in a script SQL struct.
func SQLParse(sqljson []byte) (*SQL, error) {
	var data SQL
	if err := json.Unmarshal(sqljson, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
