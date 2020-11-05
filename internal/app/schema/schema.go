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

package schema

import (
	"context"
	"time"

	"github.com/napptive/rdbms/pkg/rdbms"

	"github.com/napptive/rdbms/internal/pkg/config"
	"github.com/napptive/rdbms/internal/pkg/script"

	"github.com/jackc/pgx/v4"
)

//Load creates the basic information in the target database.
func Load(path string, defaultTimeout time.Duration, cfg config.Config) error {
	script, err := script.SQLFileParse(path)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := pgx.Connect(ctx, cfg.ConnString)
	defer conn.Close(ctx)

	if err != nil {
		return err
	}

	for _, step := range script.Steps {
		duration, err := step.TimeoutDuration(defaultTimeout)
		if err != nil {
			return err
		}

		batch := pgx.Batch{}
		for _, q := range step.Queries {
			batch.Queue(q)
		}

		bathcCtx, cancel := context.WithTimeout(ctx, duration)
		defer cancel()
		err = rdbms.ExecBatch(bathcCtx, conn, &batch)
		if err != nil {
			return err
		}
	}
	return nil
}
