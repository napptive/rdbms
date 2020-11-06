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

package rdbms

import (
	"context"

	"github.com/napptive/rdbms/internal/pkg/config"
	"github.com/napptive/rdbms/pkg/operations"
)

//Ping verifies if the database is alive.
func Ping(cfg config.Config) error {

	err := operations.Ping(context.Background(), cfg.ConnString, cfg.PingRetries, cfg.PingWaitingPeriod)
	if err != nil {
		return err
	}

	return nil
}
