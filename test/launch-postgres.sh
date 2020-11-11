#!/bin/bash

docker run -d --name local-postgres -e POSTGRES_PASSWORD=Pass2020! -p 5432:5432 postgres:13-alpine