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

package operations

import (
	"context"
	"time"

	"github.com/napptive/rdbms/v2/pkg/rdbms"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v4"
)

//ExecBatch execute a set of intructions, and stop if any instruction fails. Don't support search queries.
func ExecBatch(ctx context.Context, conn rdbms.Conn, name string, batch *pgx.Batch) error {
	result := conn.SendBatch(ctx, batch)
	defer result.Close()
	for i := 0; i < batch.Len(); i++ {
		ct, err := result.Exec()
		if err != nil {
			return err
		}
		log.Debug().Str("name", name).Int("id", i).
			Int64("rows-affected", ct.RowsAffected()).Msgf("Query (%d) succesfully executed", i)
	}
	return nil
}

// Ping execute a ping to a database n times waiting s time
func Ping(ctx context.Context, connstring string, n int, s time.Duration) error {
	r := rdbms.NewRDBMS()
	return ping(ctx, connstring, n, s, r)
}

func ping(ctx context.Context, connstring string, n int, s time.Duration, rdbms rdbms.RDBMS) error {
	var err error
	for i := 0; i < n; i++ {
		err = pingConn(ctx, connstring, s, rdbms)
		if err != nil && i != n {
			time.Sleep(s)
		}
	}
	return err
}

func pingConn(ctx context.Context, connstring string, s time.Duration, rdbms rdbms.RDBMS) error {
	timeCtx, cancel := context.WithTimeout(ctx, s)
	defer cancel()

	conn, err := rdbms.SingleConnect(timeCtx, connstring)
	if err == nil {
		defer conn.Close(timeCtx)
		return conn.Ping(timeCtx)
	}
	return err
}
