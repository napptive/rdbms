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

	"github.com/rs/zerolog/log"
)

// Config structure with all the options required by the service and service components.
type Config struct {
	RDBMS

	Version string
	Commit  string
}

// IsValid checks if the configuration options are valid.
func (c *Config) IsValid() error {
	if c.Version == "" {
		return errors.New("Version is empty")
	}
	if c.Commit == "" {
		return errors.New("Commit is empty")
	}
	err := c.RDBMS.IsValid()
	if err != nil {
		return err
	}
	return nil
}

// Print the configuration using the application logger.
func (c *Config) Print() {
	// Use logger to print the configuration
	log.Info().Str("version", c.Version).Str("commit", c.Commit).Msg("Config file")
}
