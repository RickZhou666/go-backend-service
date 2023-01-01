# go-backend-service

go-backend-service

[class udemy page](https://pplearn.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/learn/lecture/25820662?learning_path_id=4257034#overview)

[github/techschool/simplebank](https://github.com/techschool/simplebank)

| Key                              | value                                   |
| -------------------------------- | --------------------------------------- |
| `$ history \| grep "docker run"` | display all the commands                |
| offset and limit                 | https://dataschool.com/learn-sql/limit/ |
| `\* \&` sign in golang           | https://go.dev/tour/moretypes/1         |
| what is go routines              | https://go.dev/tour/concurrency/1       |

<br><br>

Postgres command<br>

| Key                                | value              |
| ---------------------------------- | ------------------ |
| `SELECT datname FROM pg_database;` | show all databases |
| `SELECT now();`                    |                    |
| `\q;`                              |                    |
| `\c simple_bank;`                  |                    |
| `\dt;`                             | show all tabls     |

# Section1: Working with database [Postgres + SQLC]

<br><br><br>

1. https://github.com/techschool/simplebank
2. https://dbdiagram.io/home

<br><br>

## 1.1 install docker + postgres + tableplus

[postgres](https://hub.docker.com/_/postgres)<br>
[tableplus](https://tableplus.com/)<br>

```bash
# 1. install postgres
$ docker pull postgres:15-alpine

# 2. display images
$ docker images

# 3. start postgres containter
$ docker run -d -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 --name postgres15 postgres:15-alpine

# 4. connect to container
# -it: run command as an interactive TTY session
# $ docker exec -it postgres15 psql -U postgres
$ docker exec -it postgres15 psql -U root

# display now
>> select now();
# quit
>> \q

# 5. display logs
$ docker logs postgres15

# 6. connect to db

```

![imgs](./imgs/Xnip2022-12-26_12-11-32.jpg)

```bash
# 7. tableplus
# (1)  execute all queries
>> command + entry

# (2) refresh
>> command + R

# 8. stop/start
$ docker stop postgres15
$ docker start postgres15

# 9. run shell in containter
$ docker exec -it postgres15 /bin/sh

```

<br><br>

## 1.2 write & run database migration in goland

```bash
# 1. install golang-migrate
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
$ brew install golang-migrate

# 2. check version
$ migrate -version


# 3. create db migration files
$ migrate create -ext sql -dir db/migration -seq init_schema

# 4. create db
$ createdb --username=root --owner=root simple_bank

# 5. acess db
$ psql simple_bank

# 6. delete db
$ dropdb simple_bank

# 7. exit
$ exit

# 8. create db from external
$ docker exec -it postgres15 createdb --username=root --owner=root simple_bank_test

# 9. access db directly
$ docker exec -it postgres15 psql -U root simple_bank_test

# 10. migrate db up
$ migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

# 11. migrate db down
$ migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
```

<br><br>

## 1.3 Generate CURD golang code from SQL

[gorm](https://gorm.io/) could be very slow when there's high load <br>
[sqlx](https://github.com/jmoiron/sqlx) <br>
[sqlc](https://github.com/kyleconroy/sqlc) <br>

```bash
# (1) install sqlc
$ brew install sqlc

# (2)
$ sqlc init

# (3) generate CRUD code
$ sqlc generate

```

### go mod init

```bash
# (1) go mod init
$ go mod init github.com/RickZhou666/go-backend-service

# (2) install go dependencies
$ go mod tidy

# (3) don't modify auto generate sqlc go files

# (4) execute only
-- name: updateAccount :exec
UPDATE accounts
SET balance = $2
WHERE id = $1;

# (5) execute with return
-- name: updateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;
```

<br><br>

## 1.4 unit test for CRUD

[golang lib pq](https://github.com/lib/pq) <br>
[golang testify](https://github.com/stretchr/testify) <br>

```bash
# -v            :for verbose log
# -cover        :to measure code coverage
# ./...         :to run all uts
# -count=1      :to disable cache. otherwise it will read from cache directly
go test -v --cover -count=1 ./...
```

if has `Testxxxx` prefix, it will run as unit test<br>
`func TestCreateaccount(t *testing.T) {`

`run package tests` - run all unit tests<br>
`run test` - run single unit test<br>
![imgs](./imgs/Xnip2022-12-31_23-02-53.jpg)

<br><br>

## 1.5 A clean way to implement database transaction in Golang

1. create a transfer record with amount = 10;
2. create an account entry for account 1 with amount = -10
3. create an account entry for account 2 with amount = +10
4. subtract 10 from the balance of account 1
5. add 10 to the balance of account 2

<br>

### 1.5.1. why do we need transaction?

1. to provide a reliable and consistent unit of work, even in case of system failure
2. to provide isolation between programs that access the database concurrently

<br>

### 1.5.2 ACID property

1. atomicity (A)<br>
   either all operations complete successfully or the transaction fails and the db is unchanged

2. consistentcy (C)<br>
   the db state must be valid after the transaction. All constraints must be satisfied

3. Isolation (I)<br>
   Concurrent transactions must not affect each other

4. Durability (D)<br>
   data written by a successful transaction must be recorded in a persistent storage

<br>

### 1.5.3 how to run SQL TX?

```go
// 1. make transaction
BEGIN;
...
COMMIT;

// 2. rollback when failure
BEGIN;
...
ROLLBACK;
```
