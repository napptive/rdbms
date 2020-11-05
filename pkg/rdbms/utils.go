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

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v4"
)

//ExecBatch execute a set of intructions, and stop if any instruction fails. Don't support search queries.
func ExecBatch(ctx context.Context, conn *pgx.Conn, batch *pgx.Batch) error {
	result := conn.SendBatch(ctx, batch)
	for i := 0; i < batch.Len(); i++ {
		ct, err := result.Exec()
		if err != nil {
			return err
		}
		log.Debug().Int("id", i).Int64("rows-affected", ct.RowsAffected()).Msgf("Query (%d) succesfully executed", i)
	}
	return nil
}
