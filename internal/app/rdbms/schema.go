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
	"time"

	"github.com/napptive/rdbms/internal/pkg/config"
	"github.com/napptive/rdbms/internal/pkg/script"
	"github.com/napptive/rdbms/pkg/operations"

	"github.com/jackc/pgx/v4"
)

//Load creates the basic information in the target database.
func Load(path string, defaultTimeout time.Duration, cfg config.Config) error {
	if !cfg.SkipPing {
		if err := Ping(cfg); err != nil {
			return err
		}
	}

	script, err := script.SQLFileParse(path)
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(context.Background(), cfg.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	for _, step := range script.Steps {
		duration, err := step.TimeoutDuration(defaultTimeout)
		if err != nil {
			return err
		}

		batch := pgx.Batch{}
		for _, q := range step.Queries {
			batch.Queue(q)
		}

		bathcCtx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()
		err = operations.ExecBatch(bathcCtx, conn, &batch)
		if err != nil {
			return err
		}
	}
	return nil
}
