# Simple transactions processing app

### Description

Small application for transaction processing.
Accepts and saves transactions in PostgresDb.
Creates single user with balance, each transaction changes balance. 
Solves concurrency problems through mutex.

Exposes single endpoint:

```
POST /tx HTTP/1.1
```

Accepts request:
```
Source-Type: client
Content-Type: application/json
{"state": "lost", "amount": "11.35", "transactionId": "unique id"} 
```

State values: win or lost.
Amount is decimal value, greater then zero. TransactionId is unique id of transaction.

Each *cancellation_interval* minutes last 10 odd transactions canceled and balance corrected.

### Configuration

Configuration should be in *config.yaml*, in the same folder with application binary.

* db: database connection credentials
* sources: accepted *Source-Type* header values, by default: game, server and payment
* cancellation_interval: cancellation interval length in minutes

### Requirements

* Docker v19+
* Docker Compose v1.24+
* Make
* Shell script support
* Golang 1.13+ (for local build only)

### Components

* Go Application
* Postgres database
* Adminer - web interface to interact with database

### Run and testing

Build and run app with PostgresDb and Adminer in docker docker containers:
```
make up
```

Transactions endpoint: http://localhost:8080/tx

Adminer database web interface http://localhost:4001

Default credentials:
```
dbname: txpost
user: postgres
password: txpost_password
```

Send test transactions using Postman collection and Newman:
```
make seed
```

Clean test network and containers:
```
make clean
```

Build local binary:
```
make build
```