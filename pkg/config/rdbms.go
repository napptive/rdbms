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

package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
)

// RDBMS contains the configuration information of the RDBMS connection. Use this configuration in the
// components requiring connections to standarize the way we ask for database information.
type RDBMS struct {
	// ConnStr is the connection string parameter.
	ConnStr string

	// Schema is the postgres schema where the data is stored.
	Schema string

	// QueryTimeout is the timeout for the query request.
	QueryTimeout time.Duration
}

// IsValid checks if the configuration options are valid.
func (r *RDBMS) IsValid() error {
	if r.ConnStr == "" {
		return nerrors.NewInvalidArgumentError("ConnStr is empty")
	}
	if r.Schema == "" {
		return nerrors.NewInvalidArgumentError("Schema is empty")
	}
	return nil
}

// MasqueradeConnStr replace the password and user by '*'
func (r *RDBMS) MasqueradeConnStr() string {
	// host=localhost user=postgres password=Pass2020! port=5432
	connStr := r.ConnStr
	fields := strings.Split(connStr, " ")
	newStr := make([]string, len(fields))
	for i, field := range fields {
		split := strings.Split(field, "=")
		if split[0] == "user" || split[0] == "password" {
			split[1] = strings.Repeat("*", len(split[1]))
		}
		newStr[i] = fmt.Sprintf("%s=%s", split[0], split[1])
	}

	return strings.Join(newStr, " ")

}

// Print the configuration using the application logger.
func (r *RDBMS) Print() {
	// Use logger to print the configuration
	log.Info().
		Str("ConnStr", r.MasqueradeConnStr()).
		Str("Schema", r.Schema).
		Dur("QueryTimeout", r.QueryTimeout).
		Msg("RDBMS config")
}
