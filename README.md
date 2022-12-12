# wallet-service

This is a service that provides a wallet for the user. It is a GRPC API that can be used to create, manage and use wallets.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installing](#installing)
- [Running the tests](#running-the-tests)
- [Running the GRPC service](#running-the-grpc-service)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

Database is Postgres and the service is written in Go.

- [Database schema](/doc/schema.sql)
- [DBML](/doc/db.dbml) 
- [DBML Visualise](https://dbdocs.io/Piyushhbhutoria/wallet)

### Prerequisites

- Go
- Protobuf
- Postgres
- BloomRPC

### Installing

A step by step series of examples that tell you how to get a development env running

```bash
brew install go
brew install protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

brew install postgresql@14
brew services start postgresql@14

brew install --cask bloomrpc
```

update `DATABASE_URL` in run and test scripts

## Running the tests

```bash
./test
```

## running-the-grpc-service

```bash
./run
```

Server is now listening on port mentioned in the `run` script
use BloomRPC to test the GRPC API
