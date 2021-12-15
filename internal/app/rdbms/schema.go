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

	"github.com/rs/zerolog/log"

	"github.com/napptive/rdbms/v2/internal/pkg/config"
	"github.com/napptive/rdbms/v2/internal/pkg/script"
	"github.com/napptive/rdbms/v2/pkg/operations"
	"github.com/napptive/rdbms/v2/pkg/rdbms"

	"github.com/jackc/pgx/v4"
)

// LoadResult store the result infomation of the load operation.
type LoadResult struct {
	ExecutedSteps []string
	SkippedSteps  []string
}

// Print write in the log the result information.
func (r *LoadResult) Print() {
	log.Info().Strs("executed", r.ExecutedSteps).
		Strs("skipped", r.SkippedSteps).
		Msgf("Load has executed %d steps and skipped %d steps", len(r.ExecutedSteps), len(r.SkippedSteps))
}

//Load creates the basic information in the target database.
func Load(path string, defaultTimeout time.Duration, selectedSteps []string, cfg config.Config) (*LoadResult, error) {
	r := rdbms.NewRDBMS()
	return load(path, defaultTimeout, selectedSteps, cfg, r)
}

func load(path string, defaultTimeout time.Duration, selectedSteps []string, cfg config.Config, rdbms rdbms.RDBMS) (*LoadResult, error) {
	var executedSteps []string
	var skippedSteps []string

	if !cfg.SkipPing {
		if err := Ping(cfg); err != nil {
			return nil, err
		}
	}

	script, err := script.SQLFileParse(path)
	if err != nil {
		return nil, err
	}

	conn, err := rdbms.SingleConnect(context.Background(), cfg.ConnString)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	for _, step := range script.Steps {
		if len(selectedSteps) == 0 || contains(selectedSteps, step.Name) {
			executedSteps = append(executedSteps, step.Name)
			err := execStep(step, defaultTimeout, conn)
			if err != nil {
				return nil, err
			}
		} else {
			skippedSteps = append(skippedSteps, step.Name)
		}
	}
	result := &LoadResult{ExecutedSteps: executedSteps, SkippedSteps: skippedSteps}
	return result, nil
}

func execStep(step script.SQLStep, defaultTimeout time.Duration, conn rdbms.Conn) error {
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
	err = operations.ExecBatch(bathcCtx, conn, step.Name, &batch)
	if err != nil {
		return err
	}
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
