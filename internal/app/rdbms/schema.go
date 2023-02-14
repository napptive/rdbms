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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/napptive/rdbms/v2/internal/pkg/config"
	"github.com/napptive/rdbms/v2/internal/pkg/script"
	"github.com/napptive/rdbms/v2/pkg/operations"
	"github.com/napptive/rdbms/v2/pkg/rdbms"
	"github.com/rs/zerolog/log"
)

// LoadResult store the result information of the loadFile operation.
type LoadResult struct {
	FileName      string
	ExecutedSteps []string
	SkippedSteps  []string
}

// Print write in the log the result information.
func (r *LoadResult) Print() {
	log.Info().Str("file", r.FileName).Strs("executed", r.ExecutedSteps).
		Strs("skipped", r.SkippedSteps).
		Msgf("Load has executed %d steps and skipped %d steps", len(r.ExecutedSteps), len(r.SkippedSteps))
}

// Load creates the basic information in the target database.
func Load(path string, fileExtension string, defaultTimeout time.Duration, selectedSteps []string, cfg config.Config) ([]*LoadResult, error) {
	r := rdbms.NewRDBMS()
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return loadDirectory(path, fileExtension, defaultTimeout, selectedSteps, cfg, r)
	}
	fileResult, err := loadFile(path, defaultTimeout, selectedSteps, cfg, r)
	if err != nil {
		return nil, err
	}

	return []*LoadResult{fileResult}, nil
}

// loadDirectory processes all files matching the extension in alphabetical order.
func loadDirectory(path string, fileExtension string, defaultTimeout time.Duration, selectedSteps []string, cfg config.Config, rdbms rdbms.RDBMS) ([]*LoadResult, error) {
	results := make([]*LoadResult, 0)
	log.Info().Str("dir", path).Msg("processing directory")
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), fileExtension) {
				fileResult, err := loadFile(path, defaultTimeout, selectedSteps, cfg, rdbms)
				if err != nil {
					log.Error().Err(err).Str("file", path).Msg("error executing file steps")
					return err
				}
				results = append(results, fileResult)
			} else {
				log.Debug().Str("file", path).Msg("skipping file")
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func loadFile(path string, defaultTimeout time.Duration, selectedSteps []string, cfg config.Config, rdbms rdbms.RDBMS) (*LoadResult, error) {
	var executedSteps []string
	var skippedSteps []string
	log.Info().Str("path", path).Msg("processing file")
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
	result := &LoadResult{FileName: path, ExecutedSteps: executedSteps, SkippedSteps: skippedSteps}
	return result, nil
}

func execStep(step script.SQLStep, defaultTimeout time.Duration, conn rdbms.Conn) error {
	log.Info().Str("step", step.Name).Msg("executing step")
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
