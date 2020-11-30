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

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Conn is the interface that wraps the pgx connection.
type Conn interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

// SingleConn is a conn create using a single thread connection.
type SingleConn interface {
	Conn

	Close(context.Context) error
	Config() *pgx.ConnConfig

	Deallocate(context.Context, string) error

	IsClosed() bool
	PgConn() *pgconn.PgConn
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
}

// PoolConn is a Conn based on connection pool.
type PoolConn interface {
	Conn

	Close()
	Config() *pgxpool.Config
	Stat() *pgxpool.Stat

	AcquireConn(context.Context) (AcquiredConn, error)
}

type pool struct {
	*pgxpool.Pool
}

func (p *pool) AcquireConn(ctx context.Context) (AcquiredConn, error) {
	ac, err := p.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &acquiredConn{ac}, nil
}

// AcquiredConn is a Conn acquired from a Pool
type AcquiredConn interface {
	Conn

	SingleConn() SingleConn
	Release()
}

type acquiredConn struct {
	*pgxpool.Conn
}

func (c *acquiredConn) SingleConn() SingleConn {
	return c.Conn.Conn()
}
