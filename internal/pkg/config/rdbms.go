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
	"errors"
	"time"
)

// RDBMS is a structure with all the options required by the service to config the database connection.
type RDBMS struct {
	// ConnString is the postgres connection string.
	ConnString string

	// PingRetries is the number of retries that a ping is going to be launched.
	PingRetries int

	// PingWaitingPerios is the time that the thread wait between each ping.
	PingWaitingPeriod time.Duration

	// IgnorePing allows to the rest of commands skip the ping step
	SkipPing bool
}

// IsValid checks if the RBDMS configuration options are valid.
func (c *RDBMS) IsValid() error {
	if c.ConnString == "" {
		return errors.New("ConnString is empty")
	}
	if c.PingRetries <= 0 {
		return errors.New("PingRetries is <=0")
	}

	return nil
}
