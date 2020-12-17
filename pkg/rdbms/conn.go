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

// CommandTag is the result of an Exec function
type CommandTag interface {
	// Delete is true if the command tag starts with "DELETE".
	Delete() bool

	// Insert is true if the command tag starts with "INSERT".
	Insert() bool

	// RowsAffected returns the number of rows affected. If the CommandTag was not for a row affecting command (e.g. "CREATE TABLE") then it returns 0.
	RowsAffected() int64

	// Select is true if the command tag starts with "SELECT".
	Select() bool

	// String transform the result into a string.
	String() string

	// Update is true if the command tag starts with "UPDATE".
	Update() bool
}

// Conn is the interface that wraps the pgx connection.
type Conn interface {
	// Begin starts a transaction. Unlike database/sql, the context only affects the begin command. i.e. there is no auto-rollback on context cancellation.
	Begin(context.Context) (pgx.Tx, error)

	// BeginTx starts a transaction with txOptions determining the transaction mode. Unlike database/sql, the context only affects the begin command. i.e. there is no auto-rollback on context cancellation.
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)

	// CopyFrom uses the PostgreSQL copy protocol to perform bulk data insertion. It returns the number of rows copied and an error.
	// CopyFrom requires all values use the binary format. Almost all types implemented by pgx use the binary format by default. Types implementing Encoder can only be used if they encode to the binary format.
	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)

	// Exec executes sql. sql can be either a prepared statement name or an SQL string. arguments should be referenced positionally from the sql string as $1, $2, etc.
	Exec(context.Context, string, ...interface{}) (CommandTag, error)

	// Query executes sql with args. If there is an error the returned Rows will be returned in an error state. So it is allowed to ignore the error returned from Query and handle it in Rows.
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely needed. See the documentation for those types for details.
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while querying is deferred until calling Scan on the returned Row. That Row will error with ErrNoRows if no rows are returned.
	QueryRow(context.Context, string, ...interface{}) pgx.Row

	// SendBatch sends all queued queries to the server at once. All queries are run in an implicit transaction unless explicit transaction control statements are executed. The returned BatchResults must be closed before the connection is used again.
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

// SingleConn is a conn create using a single thread connection.
type SingleConn interface {
	// Conn is the basic pgx connection
	Conn

	// Close closes a connection. It is safe to call Close on a already closed connection.
	Close(context.Context) error

	// Config returns a copy of config that was used to establish this connection.
	Config() *pgx.ConnConfig

	// Deallocate released a prepared statement
	Deallocate(context.Context, string) error

	// IsClosed check if the connection state.
	IsClosed() bool

	// PgConn returns the underlying *pgconn.PgConn. This is an escape hatch method that allows lower level access to the PostgreSQL connection than pgx exposes.
	PgConn() *pgconn.PgConn

	// Ping send a ";" query to the connected database.
	Ping(context.Context) error

	// Prepare creates a prepared statement with name and sql. sql can contain placeholders for bound parameters. These placeholders are referenced positional as $1, $2, etc.
	// Prepare is idempotent; i.e. it is safe to call Prepare multiple times with the same name and sql arguments. This allows a code path to Prepare and Query/Exec without concern for if the statement has already been prepared.
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
}

type singleConn struct {
	*pgx.Conn
}

// Exec executes sql. sql can be either a prepared statement name or an SQL string. arguments should be referenced positionally from the sql string as $1, $2, etc.
func (c *singleConn) Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error) {
	return c.Conn.Exec(ctx, query, args)
}

// PoolConn is a Conn based on connection pool.
type PoolConn interface {
	// Conn is the basic pgx connection
	Conn
	// Close closes all connections in the pool and rejects future Acquire calls. Blocks until all connections are returned to pool and closed.
	Close()
	// Config returns a copy of config that was used to initialize this pool.
	Config() *pgxpool.Config

	// Stat return a object with pool stats.
	Stat() *pgxpool.Stat

	// AcquireConn obtains a free connection.
	AcquireConn(context.Context) (AcquiredConn, error)
}

// pool is a struct to wrapp the pgxpool.Pool struct and return rdbms interface objects.
type pool struct {
	*pgxpool.Pool
}

// Exec executes sql. sql can be either a prepared statement name or an SQL string. arguments should be referenced positionally from the sql string as $1, $2, etc.
func (p *pool) Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error) {
	return p.Pool.Exec(ctx, query, args)
}

// AcquireConn obtains a free connection.
func (p *pool) AcquireConn(ctx context.Context) (AcquiredConn, error) {
	ac, err := p.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &acquiredConn{ac}, nil
}

// AcquiredConn is a Conn acquired from a Pool
type AcquiredConn interface {
	// Conn is the basic pgx connection
	Conn

	// SingleConn recovers the underlying single conn used by the acquired connection
	SingleConn() SingleConn

	//Release returns the connection to the pool it was acquired from. Once Release has been called, other methods must not be called. However, it is safe to call Release multiple times. Subsequent calls after the first will be ignored.
	Release()
}

// acquiredConn is a struct to wrapp the pgxpool.Conn struct and return rdbms interface objects.
type acquiredConn struct {
	*pgxpool.Conn
}

// Exec executes sql. sql can be either a prepared statement name or an SQL string. arguments should be referenced positionally from the sql string as $1, $2, etc.
func (c *acquiredConn) Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error) {
	return c.Conn.Exec(ctx, query, args)
}

// SingleConn recovers the underlying single conn used by the acquired connection
func (c *acquiredConn) SingleConn() SingleConn {
	return &singleConn{c.Conn.Conn()}
}
