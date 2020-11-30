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

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//RDBMS is the interface to create rdbms connections.
type RDBMS interface {
	// SingleConnect establishes a connection with a PostgreSQL server with a connection string. See pgconn.Connect for details.
	SingleConnect(context.Context, string) (SingleConn, error)
	// SingleConnectConfig establishes a connection with a PostgreSQL server with a configuration struct. connConfig must have been created by ParseConfig.
	SingleConnectConfig(context.Context, *pgx.ConnConfig) (SingleConn, error)

	// PoolConnect creates a new Pool and immediately establishes one connection. ctx can be used to cancel this initial connection. See ParseConfig for information on connString format.
	PoolConnect(context.Context, string) (PoolConn, error)

	// PoolConnectConfig creates a new Pool and immediately establishes one connection. ctx can be used to cancel this initial connection. config must have been created by ParseConfig.
	PoolConnectConfig(context.Context, *pgxpool.Config) (PoolConn, error)
}

//NewRDBMS create a new RDBMS instance.
func NewRDBMS() RDBMS {
	return &rdbms{}
}

//rdbms is the internal struct that contains the connection methods.
type rdbms struct{}

// SingleConnect establishes a connection with a PostgreSQL server with a connection string.
func (r *rdbms) SingleConnect(ctx context.Context, connString string) (SingleConn, error) {
	return pgx.Connect(ctx, connString)
}

// SingleConnectConfig establishes a connection with a PostgreSQL server with a configuration struct. connConfig must have been created by ParseConfig.
func (r *rdbms) SingleConnectConfig(ctx context.Context, connConfig *pgx.ConnConfig) (SingleConn, error) {
	return pgx.ConnectConfig(ctx, connConfig)
}

// PoolConnect creates a new Pool and immediately establishes one connection. ctx can be used to cancel this initial connection. See ParseConfig for information on connString format.
func (r *rdbms) PoolConnect(ctx context.Context, connString string) (PoolConn, error) {
	p, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &pool{p}, nil
}

// PoolConnectConfig creates a new Pool and immediately establishes one connection. ctx can be used to cancel this initial connection. config must have been created by ParseConfig.
func (r *rdbms) PoolConnectConfig(ctx context.Context, config *pgxpool.Config) (PoolConn, error) {
	p, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return &pool{p}, nil
}
