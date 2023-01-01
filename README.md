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

| Key                                                                          | value              |
| ---------------------------------------------------------------------------- | ------------------ |
| `SELECT datname FROM pg_database;`                                           | show all databases |
| `SELECT now();`                                                              |                    |
| `\q;`                                                                        |                    |
| `\c simple_bank;`                                                            |                    |
| `\dt;`                                                                       | show all tables    |
| INSERT INTO accounts (owner, balance, currency) VALUES ('rick', 100, 'USD'); | insert column      |

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
$ docker exec -it postgres15 createdb --username=root --owner=root simple_bank

# 9. access db directly
$ docker exec -it postgres15 psql -U root simple_bank

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

facing github link issue, remove from local and clone from remote<br>
![imgs](./imgs/Xnip2023-01-01_14-22-13.jpg)

<br><br>

## 1.6 DB transaction lock & how to handle deadlock in Golang

[postgresql lock monitoring](https://wiki.postgresql.org/wiki/Lock_Monitoring)

<br><br>

### 1.6.1 test in terminal for 2 concurrent txn

```bash
# 1. start txn for 1st
>> BEGIN;
>> select * from accounts where id = 126 for update;

# 2. start txn for 2nd
>> BEGIN;
>> select * from accounts where id = 126 for update;
# the 2nd will get blocked

# 3. update for the 1st
>> update accounts set balance = 500 where id = 126;
>> COMMIT;

# 4. after we commit 1st, the 2nd is unblocked with latest result

```

### 1.6.2 update get query

```sql
-- (1) add get for update query
-- name: getAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- (2) update query
$ make sqlc
```

### 1.6.3 deadlock detected

<br><br>

#### 1.6.3.1 replicate in terminal

<br><br>

#### 1.6.3.2 way1: remove foreign key

```bash
# (1) comment foregin key line in schema_up
# ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
# ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
# ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

# (2) delete db
$ make migratedown

# (3) regenerate db
$ make migrateup
```

<br><br>

#### 1.6.3.3 way2: remove foreign key binding in get query

```bash
# (1) as the foreign key is linked to id which is primary key, and we won't change it
#     so we add NO KEY in get query
-- name: getAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

# (2) regenerate
$ make sqlc


```

<br><br>

## 1.7 how to avoid db deadlock?

### 1.7.1 replicate deadlock

```sql
-- 1. run below in 1st
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;

-- 2. run below in 2nd
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;

-- 3. then run below in 1st
UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;
-- then it's getting stuck
-- as txn 2 tried to get a sharelock 908. but 908 is exclusivelock and holding by txn 1
```

![imgs](./imgs/Xnip2023-01-01_23-00-51.jpg)

```sql
-- 4. if we keeping running below in 2nd
UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
-- we'll observed deadlock error

```

![imgs](./imgs/Xnip2023-01-01_23-03-19.jpg)

<br><br>

### 1.7.2 avoid deadlock

`execute account in same order`

```sql
-- 1. run below in 1st
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;

-- 2. run below in 2nd
BEGIN;

UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
-- 2nd will be block

-- 3. run below in 1st
UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;

COMMIT;
-- 2nd will return value
```

`update account with smaller id always`
<br><br>

## 1.8 Deeply understand txn isolation levels & read phenomena

### 1.8.1 Read Phenomena

1. Diry read<br>
   A txn `reads` data writtern by other concurrent `uncommited` txn<br>

2. Non-repeatable read<br>
   A txn `reads` the `same row twice` and sees different value because it has been `modified` by other `commited` txn<br>

3. Phantom read<br>
   A txn `re-executes` a query to `find rows` that satisfy a condition and sees a `different set` of rows, due to changes by other `commited` txn<br>

4. Serilization anomaly<br>
   The result of a `group` of concurrent `commited txns` is `impossible to achieve` if we try to run them `sequentially` in any order without overlapping<br>

<br><br>

### 1.8.2 4 Standard Isolation Levels

American National Standards Insititues - ANSI
![imgs](./imgs/Xnip2023-01-01_23-57-49.jpg)

<br><br>

### 1.8.3 Mysql

```bash
# (1) pull mysql
$ docker pull mysql:

# (2) run image in container
$ docker run -d -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 --name postgres15 postgres:15-alpine

# (3) exec containter
$ docker exec -it mysql8 mysql -uroot -psecret simple_bank

# (4) check txn isolation
>> select @@transaction_isolation;
>> select @@global.transaction_isolation;

# (5) change txn isolation level
>> set session transaction isolation level read uncommited;
```

<br><br>

#### 1.8.3.0 Isolation levels in MySQL

![imgs](./imgs/Xnip2023-01-02_01-06-08.jpg)

<br><br>

#### 1.8.3.1 read uncommited

<br><br>

#### 1.8.3.2 read commited

<br><br>

#### 1.8.3.3 repeatable read

will read same row, but result is not consistent. If another txn increase by 10 to 110, current txn will still read it as 100, if we increase by 10 as well in conrrent txn, it will change to 120. but it's from 100 -> 120

<br><br>

#### 1.8.3.4 serializable

1. the read txn will block any update or delete txn in same row
2. if update in both txn, deadlock will occur
3. if we execute both select in 2 txn, it's ok. then we execute insert in txn1, it will block as txn2 is holding a share lock. after we execute insert in txn2, it will fail and release lock, txn1 can execute.

<br><br>

### 1.8.4 postgres

```sql
-- (1) login to psql
$ docker exec -it postgres15 psql -U root
-- (2) check txn isolation level
>> show transaction isolation level;
-- default as read committed
```

<br><br>

#### 1.8.4.0 Isolation levels in Postgres

![imgs](./imgs/Xnip2023-01-02_01-06-55.jpg)

<br><br>

#### 1.8.4.1 read uncommited

```sql
-- (1) postgres can only set txn isolation level during txn
>> BEGIN;
>> set transaction isolation level read uncommitted;
>> show transaction isolation level;

-- (2) select all accounts
>> SELECT * FROM accounts;

-- (3) checking id 1 in txn 2
>> SELECT * FROM accounts where id = 1;
-- result in $100

-- (4) update balance in txn 1
>> update accounts set balance = balance - 10 where id = 1 returning *;
-- result in $90

-- (5) checking id 1 in txn 2
>> SELECT * FROM accounts where id = 1;
-- result still in $100

-- (6) commit in txn 1
>> COMMIT;

-- (7) checking id 1 in txn 2
>> SELECT * FROM accounts where id = 1;
-- result in $90
```

[PostgreSQL's Read Uncommitted mode behaves like Read Committed](https://www.postgresql.org/docs/9.5/transaction-iso.html)

<br><br>

#### 1.8.4.2 read commited

`phantom read occured`

```sql
-- (1) postgres can only set txn isolation level during txn
>> BEGIN;
>> set transaction isolation level read committed;
>> show transaction isolation level;

-- (2) select all accounts in txn 1
>> SELECT * FROM accounts;

-- (3) checking id 1 in txn 2
>> SELECT * FROM accounts where id = 1;
-- result in $90

-- (4) checking all balance >= 90 in txn 2
>> SELECT * FROM accounts where balance >= 90;
-- result in 2 records

-- (5) update balance in txn 1
>> update accounts set balance = balance - 10 where id = 1 returning *;
-- result in $80

-- (6) commit in txn 1
>> COMMIT;

-- (7) checking id 1 in txn 2
>> SELECT * FROM accounts where id = 1;
-- result will be in $80

-- (8) checking all balance >= 90 in txn 2
>> SELECT * FROM accounts where balance >= 90;
-- we will only see 1 record, so Phantom read occured

-- (9) commit txn 2
>> commit;
```

<br><br>

#### 1.8.4.3 repeatable read

1. `phantom read` is prevented
2. prevent update on row that is changed in another txn.
3. `serialization anomaly` occured

```sql
-- (1) postgres can only set txn isolation level during txn
>> BEGIN;
>> set transaction isolation level repeatable read;
>> show transaction isolation level;

-- (2) select all accounts in txn 1
>> SELECT * FROM accounts;

-- (3) checking id 1 in txn 2
>> SELECT * FROM accounts where id = 1;
-- result in $90

-- (4) checking all balance >= 80 in txn 2
>> SELECT * FROM accounts where balance >= 80;
-- result in 2 records

-- (5) update balance in txn 1
>> update accounts set balance = balance - 10 where id = 1 returning *;
-- result in $70

-- (6) commit in txn 1
>> COMMIT;

-- (7) checking id 1 in txn 2
>> SELECT * FROM accounts where id = 1;
-- result will still be in $80

-- (8) checking all balance >= 80 in txn 2
>> SELECT * FROM accounts where balance >= 80;
-- we will only see 2 records, so Phantom read is prevented

-- (9) update id 1 in txn 2
>> update accounts set balance = balance - 10 where id = 1 returning *;
-- ERROR:  current transaction is aborted, commands ignored until end of transaction block

-- (9) commit txn 2
>> commit;
```

<br><br>

`replicate serialization anomaly`

```sql
-- (1) postgres can only set txn isolation level during txn
>> BEGIN;
>> set transaction isolation level repeatable read;
>> show transaction isolation level;

-- (2) select all accounts in txn 1
>> SELECT * FROM accounts;

-- (3) sum balance and assign to a new account in txn 1
>> SELECT sum(balance) FROM accounts;
>> INSERT INTO accounts(owner, balance, currency) VALUES ('sum', 170, 'USD');
>> SELECT * FROM accounts;
-- result in 3 records

-- (4) checking all in txn 2
>> SELECT * FROM accounts;
-- result still in 2 records. as we use repeatable read level


-- (5) sum balance and assign to a new account in txn 2
>> SELECT sum(balance) FROM accounts;
>> INSERT INTO accounts(owner, balance, currency) VALUES ('sum', 170, 'USD');
>> SELECT * FROM accounts;
-- result in 3 records

-- (6) commit both txn 1 and txn 2
>> COMMIT;

-- (7) check all in txn 2
>> SELECT * FROM accounts;
-- there are two sum records. this is serialization anomaly
```

#### 1.8.4.4 serialization

1. `serialization anomaly` is prevented

```sql
-- (1) postgres can only set txn isolation level during txn
>> BEGIN;
>> set transaction isolation level serializable;
>> show transaction isolation level;

-- (2) select all accounts in txn 1
>> SELECT * FROM accounts;

-- (3) sum balance and assign to a new account in txn 1
>> SELECT sum(balance) FROM accounts;
>> INSERT INTO accounts(owner, balance, currency) VALUES ('sum', 510, 'USD');
>> SELECT * FROM accounts;
-- result in 5 records with new added record

-- (4) checking all in txn 2
>> SELECT * FROM accounts;
-- result still in 2 records. as we use repeatable read level


-- (5) sum balance and assign to a new account in txn 2
>> SELECT sum(balance) FROM accounts;
>> INSERT INTO accounts(owner, balance, currency) VALUES ('sum', 510, 'USD');
>> SELECT * FROM accounts;
-- result is identical with txn 1

-- (6) commit both txn 1 and txn 2
>> COMMIT;
-- txn 1 is ok

>> COMMIT;
-- tnx 2 is failed
-- ERROR:  could not serialize access due to read/write dependencies among transactions
-- DETAIL:  Reason code: Canceled on identification as a pivot, during commit attempt.
-- HINT:  The transaction might succeed if retried.

-- (7) check all in txn 2
>> SELECT * FROM accounts;
-- records stay the same. serialization anomaly is prevented
```

<br><br>

### 1.8.5 Compare MySQL vs Postgres

| item             | MySQL             | PostgreSQL             |
| ---------------- | ----------------- | ---------------------- |
| isolation levels | 4                 | 3                      |
| mechanism        | locking mechanism | dependencies detection |
| default level    | repeatable read   | read committed         |

`Keep in mind`

1. retry mechanism <br>
   there migh be errors, timeout or deadlock
2. read documentation <br>
   each database engine might implement isolation level differently

<br><br>

## 1.9 setup github actions for golang + postgres to run automated tests

<br><br><br>

# Section2: Building RESTful HTTP JSON API [Gin + JWT + PASETO]

<br><br><br>

# Section3: Deploying the application to production [Docker + Kubernetes + AWS]

<br><br><br>

# Section4: Advanced backend Topics [gRPC]

<br><br><br>
