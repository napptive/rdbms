# RDBMS
RDBMS is an application that contains utils and helper for relationa databases. The objective of this project is to speed up the integration with this type of databases.

## Project components

The project compose of the following components:
* **Postgres K8S deployment**: The K8S yaml file to deploy a postgres instance in Kubernetes.
* **RDBMS CLI**: A command line interface to load data and ping an existing database.
* **RDBMS Job Example**: An example to launch a K8S job to load data in a database.
* **RDBMS Conn Wrapper**: This project include interfaces to create PGX Connections that supports mockgen (see pkg/rdbms).

## Usage

```
This command contail useful operations to manage a Postgress database.

Usage:
  rdbms [flags]
  rdbms [command]

Examples:

  $ rdbms help
  $ rdbms ping -c "host=localhost user=postgres password=Pass2020! port=5432"
  $ rdbms schema load --scriptLoadPath test/data/ValidSQLScript.yaml
  $ rdbms schema load --scriptLoadPath test/data/ValidSQLScript.yaml --selectedStep creation-step --selectedStep drop-step 


Available Commands:
  help        Help about any command
  ping        Ping databases
  schema      Commands related with the database schema

Flags:
  -c, --connectionString string      Database connection string (default "host=localhost user=postgres password=Pass2020! port=5432")
  -l, --consoleLogging               Pretty print logging
  -d, --debug                        Set debug level
  -h, --help                         help for rdbms
  -r, --pingRetries int              Number of retries to ping to the database. (default 3)
  -w, --pingWaitingPeriod duration   Waiting time between each ping command (default 5s)
  -k, --skipPing                     If true, the command skip the ping step.

Use "rdbms [command] --help" for more information about a command.
```


## Badges

![Check changes in the Main branch](https://github.com/napptive/rdbms/workflows/Check%20changes%20in%20the%20Main%20branch/badge.svg?branch=main)


## License

 Copyright 2020 Napptive

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
