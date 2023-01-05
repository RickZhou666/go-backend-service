# go-backend-service

go-backend-service

<br><br>

## 0.1 Questions

1. when to use uppercase for func? when to use lowercase for func?

```go
func IsSupportedCurrency(currency string) bool {}

func errorResponse(err error) gin.H {}
```

2. Difference between `*` and `&`

3. what is BTREE?

<br><br>

## 0.2 Tips

[class udemy page](https://pplearn.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/learn/lecture/25820662?learning_path_id=4257034#overview)

[github/techschool/simplebank](https://github.com/techschool/simplebank)

| Key                              | value                                                                                       |
| -------------------------------- | ------------------------------------------------------------------------------------------- |
| `$ history \| grep "docker run"` | display all the commands                                                                    |
| offset and limit                 | https://dataschool.com/learn-sql/limit/                                                     |
| `\* \&` sign in golang           | https://go.dev/tour/moretypes/1                                                             |
| what is go routines              | https://go.dev/tour/concurrency/1                                                           |
| `:=` and `=`                     | `:=` 初始化并赋值，回覆盖原来的值 <br> `=` 直接赋值                                         |
|                                  | https://blog.csdn.net/Ivan45007/article/details/121978869                                   |
|                                  | [Short variable declarations](https://go.dev/ref/spec#Short_variable_declarations)          |
|                                  | a short variable declaration may redeclare variables provided they were originally declared |
| `curl` man page                  | https://curl.se/docs/manpage.html                                                           |
| dbdiagram                        | https://dbdiagram.io                                                                        |

<br><br>

Postgres command<br>

| Key                                                                          | value                                                      |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------- |
| `SELECT datname FROM pg_database;`                                           | show all databases                                         |
| `SELECT now();`                                                              |                                                            |
| `\q;`                                                                        |                                                            |
| `\c simple_bank;`                                                            |                                                            |
| `\dt;`                                                                       | show all tables                                            |
| INSERT INTO accounts (owner, balance, currency) VALUES ('rick', 100, 'USD'); | insert column                                              |
| `LIMIT` and `OFFSET`                                                         | https://www.postgresql.org/docs/current/queries-limit.html |
| `LIMIT` and `OFFSET`                                                         | `LIMIT=5` at most 5 records in one page                    |
| `LIMIT` and `OFFSET`                                                         | `OFFSET=10` skip first 10 records                          |

<br><br>

TablePlus command<br>

| Key          | Value                           |
| ------------ | ------------------------------- |
| refresh data | command + R                     |
| delete data  | 1. delete 2. command + S        |
| modify data  | 1. do the change 2. command + S |

<br><br><br>

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

[create table](https://www.w3schools.com/mysql/mysql_create_table.asp)

```bash
# (1) pull mysql
$ docker pull mysql:8

# (2) run image in container
$ docker run -e MYSQL_ROOT_PASSWORD=secret -p 3310:3306 -d --name mysql8 mysql:8

# (4) create database
$ show databases;
$ create database simple_bank;
$ use simple_bank;

# (3) exec containter
$ docker exec -it mysql8 mysql -uroot -psecret simple_bank



# (4) check txn isolation
>> select @@transaction_isolation;
>> select @@global.transaction_isolation;
```

<br><br>

#### 1.8.3.0 Isolation levels in MySQL

![imgs](./imgs/Xnip2023-01-02_01-06-08.jpg)

<br><br>

#### 1.8.3.1 read uncommited

1. facing `dirty read`

```sql
-- (1) change txn isolation level to read uncommitted in both txns
>> set session transaction isolation level read uncommitted;
>> select @@transaction_isolation;

-- (2) start txn
>> begin;
-- or
>> start transaction;

-- (3) select all from txn 1
>> select * from accounts;

-- (4) select id 1 from txn 2
>> select * from accounts where id = 1;

-- (5) decrease $10 from txn 1 for id 1
>> update accounts set balance = balance - 10 where id = 1;
>> select * from accounts where id = 1;
-- result in $90 for id 1

-- (6) select id 1 from txn 2
>> select * from accounts where id = 1;
-- it display the modifed uncommited balance for id 1

```

<br><br>

#### 1.8.3.2 read commited

1. avoid `dirty read`
2. facing `non-repetable read`
3. facing `phantom read`

```sql
-- (1) change txn isolation level to read committed in both txns
>> set session transaction isolation level read committed;
>> select @@transaction_isolation;

-- (2) start txn
>> begin;
-- or
>> start transaction;

-- (3) select all from txn 1
>> select * from accounts;

-- (4) select id 1 from txn 2
>> select * from accounts where id = 1;

-- (5) decrease $10 from txn 1 for id 1
>> update accounts set balance = balance - 10 where id = 1;
>> select * from accounts where id = 1;
-- result in $80 for id 1

-- (6) select id 1 from txn 2
>> select * from accounts where id = 1;
-- result still in $90. avoid dirty read for uncommited txn

-- (7) select all from txn 2
>> select * from accounts where balance >= 90;
-- 2 records

-- (8) commit txn 1
>> commit;

-- (9) check id 1 in txn 2
>> select * from accounts where id = 1;
-- result in $80 for id 1. non-repeatable read

-- (10) select all from txn 2
>> select * from accounts where balance >= 90;
-- only 1 records. facing phantom read

-- (11) commit
>> commit;
```

<br><br>

#### 1.8.3.3 repeatable read

1. avoid `dirty read`
2. avoid `non-repeatable read`
3. avoid `phantom read`
4. the result is not `consistent`

will read same row, but result is not consistent. If another txn increase by 10 to 110, current txn will still read it as 100, if we increase by 10 as well in conrrent txn, it will change to 120. but it's from 100 -> 120

```sql
-- (1) change txn isolation level to read committed in both txns
>> set session transaction isolation level repeatable read;
>> select @@transaction_isolation;

-- (2) start txn
>> begin;
-- or
>> start transaction;

-- (3) select all from txn 1
>> select * from accounts;

-- (4) select id 1 and all records from txn 2
>> select * from accounts where id = 1;
>> select * from accounts where balance >= 80;

-- (5) decrease $10 from txn 1 for id 1
>> update accounts set balance = balance - 10 where id = 1;
>> select * from accounts where id = 1;
-- result in $70 for id 1
>> select * from accounts;

-- (6) commit txn 1
>> commit;

-- (7) select id 1 from txn 2
>> select * from accounts where id = 1;
-- result still in $80. still the old version of id 1, avoid non-repeatable read

-- (8) select all from txn 2
>> select * from accounts where balance >= 80;
-- still return 2 records. the query is repeatable

-- (9) decrease $10 from txn 2 for id 1
>> update accounts set balance = balance - 10 where id = 1;
>> select * from accounts where id = 1;
-- result in $60, which is correct value

-- (8) rollback in txn 2;
>> rollback;
```

<br><br>

#### 1.8.3.4 serializable

1. avoid `serialization anomaly`

1. the read txn will block any update or delete txn in same row
1. if update in both txn, deadlock will occur
1. if we execute both select in 2 txn, it's ok. then we execute insert in txn1, it will block as txn2 is holding a share lock. after we execute insert in txn2, it will fail and release lock, txn1 can execute.

```sql
-- (1) change txn isolation level to read committed in both txns
>> set session transaction isolation level serializable;
>> select @@transaction_isolation;

-- (2) start txn
>> begin;
-- or
>> start transaction;

-- (3) select all from txn 1
>> select * from accounts;

-- (4) select id 1 and all records from txn 2
>> select * from accounts where id = 1;

-- (5) decrease $10 from txn 1 for id 1
>> update accounts set balance = balance - 10 where id = 1;
-- the update in txn 1 is blocked. In serializable level, mysql convert all plain SELECT query to SELECT FOR SHARE. the txn holding SELECT FOR SHARE only allow other txn to READ not UPDATE or DELETE.
-- ERROR 1205 (HY000): Lock wait timeout exceeded; try restarting transaction
-- make sure u use retry mechanism when u use serializble isolation level

-- (6) restart txn in txn 1
>> select @@transaction_isolation;
>> begin;
>> select * from accounts where id = 1;
>> update accounts set balance = balance - 10 where id = 1;
-- txn 1 blocked

-- (7) update in txn 2
>> update accounts set balance = balance - 10 where id = 1;
-- now txn 2 is also waiting for share lock from txn 1
-- ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction
-- txn 1 successful

-- (8) rollback both txns
>> rollback;

-- (9) restart txns
>> select @@transaction_isolation;
>> begin;

-- (10) select id 1 from txn 1
>> select * from accounts where id = 1;

-- (11) select id 1 from txn 2
>> select * from accounts where id = 1;

-- (12) update in txn 1
>> update accounts set balance = balance - 10 where id = 1;
-- txn 1 blocked

-- (13) commit txn 2
>> commit;
-- share lock released, txn 1 updated
-- account balance updated in txn 1

-- (14) commit txn 1
>> commit;
```

<br><br>

#### 1.8.3.5 how mysql handle serialization anomaly

1. prevent `serialization anomaly`

```sql
-- (1) change txn isolation level to read committed in both txns
>> set session transaction isolation level serializable;
>> select @@transaction_isolation;

-- (2) start txn
>> begin;
-- or
>> start transaction;

-- (3) select all from txn 1
>> select * from accounts;

-- (4) compute sum in txn 1
>> select sum(balance) from accounts;

-- (5) insert sum in txn 1
>> insert into accounts (owner, balance, currency) values ('sum', '160', 'usd');
>> select * from accounts;


-- (6) select all in txn 2
>> select * from accounts;
-- txn 2 is blocked, need txn 1 finish the txn

-- (7) commit in txn 1
>> commit;
-- txn 2 display the select result

-- (8) insert sum from txn 2
>> select sum(balance) from accounts;
>> insert into accounts (owner, balance, currency) values ('sum', '320', 'usd');
>> select * from accounts;
>> commit;
-- prevent serialization anomaly

```

`try different orders for both txns`

```sql
-- (1) change txn isolation level to read committed in both txns
>> set session transaction isolation level serializable;
>> select @@transaction_isolation;

-- (2) start txn
>> begin;
-- or
>> start transaction;

-- (3) select all and compute sum from txn 1
>> select * from accounts;
>> select sum(balance) from accounts;

-- (4) select all and compute sum from txn 2
>> select * from accounts;
>> select sum(balance) from accounts;
-- same sum as txn 1

-- (5) insert sum in txn 1
>> insert into accounts (owner, balance, currency) values ('sum', '640', 'usd');
-- txn 1 is blocked

-- (5) insert sum in txn 2
>> insert into accounts (owner, balance, currency) values ('sum', '640', 'usd');
-- ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction
-- txn 2 deadlock and txn 1 is successful

-- (6) rollback txn 2
rollback;

-- (7) select all in txn 1
>> select * from accounts;
-- new sum inserted successfully
```

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

### 1.9.1 Github Actions

github similar to below:

1. Jenkins
2. Travis
3. CircleCI

#### 1.9.1.1 Workflow

<br><br>

#### 1.9.1.2 Runner

<br><br>

#### 1.9.1.3 Job

<br><br>

#### 1.9.1.4 Step & Action

<br><br>

#### 1.9.1.5 Summary

<br><br>

### 1.9.2 Setup github Actions

1. go to github repo [go-backend-service](https://github.com/RickZhou666/go-backend-service)
2. Actions
3. Go
4. create local file first

```bash
$ mkdir -p .github/workflows
$ cd .github/workflows
$ touch ci.yml
```

5. after setup ci.yml and push to github
6. check `Actions` and ci-test
7. click init CI workflow
8. click Jobs `Test` to check status

![imgs](./imgs/Xnip2023-01-02_16-19-29.jpg)

9. declare github postgres actions [link](https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers)

10. check official docker postgres [setup](https://hub.docker.com/_/postgres)

11. setup db schema

```yml
# add one more step in yml
- name: Run migrations
  run: make migrateup
```

12. migrate cmd is not installed [golang/migrate doc](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
    ![imgs](./imgs/Xnip2023-01-02_16-27-32.jpg)

```bash
# 1. download
$ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
$ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz

# 2. move to below folder
$ mv migrate /usr/bin/
```

`copy` correct linux file
![imgs](./imgs/Xnip2023-01-02_16-30-36.jpg)

13. Run migrations still failed
    as we haven't expost the ports
    ![imgs](./imgs/Xnip2023-01-02_16-48-03.jpg)

14. CI job passed
    ![imgs](./imgs/Xnip2023-01-02_17-20-34.jpg)

<br><br><br>

# Section2: Building RESTful HTTP JSON API [Gin + JWT + PASETO]

| Key                              | Value                                                       |
| -------------------------------- | ----------------------------------------------------------- |
| Gin Model binding and validation | https://gin-gonic.com/docs/examples/binding-and-validation/ |
| sqlc config                      | https://docs.sqlc.dev/en/latest/reference/config.html       |
| viper                            | https://github.com/spf13/viper                              |

<br><br>

## 2.1 Implement RESTful HTTP API in Go using Gin

| Popular web frameworks | Popular HTTP routers |
| ---------------------- | -------------------- |
| Gin                    | FastHttp             |
| Beego                  | Gorrila Mux          |
| Echo                   | HttpRouter           |
| Revel                  | Chi                  |
| Martini                |                      |
| Fiber                  |                      |
| Buffalo                |                      |

```bash
# (1) install Gin
# https://github.com/gin-gonic/gin
$ go get -u github.com/gin-gonic/gin
$ go mod tidy

# (2) handler function declare as *Context
```

![imgs](./imgs/Xnip2023-01-03_12-40-30.jpg)

```bash
# (3) validate requests
https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Baked_In_Validators_and_Tags
# One of

# (4) get dependencies
$ go get github.com/go-delve/delve/service/api

```

<br><br>

### 2.1.1 create account

```bash
# (1) test in postman
```

![imgs](./imgs/Xnip2023-01-05_10-54-52.jpg)

```bash
# (2) displayed the response result
```

![imgs](./imgs/Xnip2023-01-05_10-55-23.jpg)

### 2.1.2 get account

```bash
# (1) test in postman
```

![imgs](./imgs/Xnip2023-01-05_11-08-10.jpg)

### 2.1.3 get list of accounts

<br><br>

`pagenation`

[emit_empty_slices](https://docs.sqlc.dev/en/latest/reference/config.html)<br>
If true, slices returned by :many queries will be empty instead of nil. Defaults to false.

```bash
# (1) upgrade sqlc to latest version
$ brew upgrade sqlc

# (2) check sqlc version
$ sqlc version

# (3) regenerate code
$ make sqlc

# (4) now items are initialized as empty slice
```

![imgs](./imgs/Xnip2023-01-05_11-27-53.jpg)

```bash
# (5) now we returned as empty array instead of null
```

![imgs](./imgs/Xnip2023-01-05_11-28-50.jpg)

<br><br>

## 2.2 load config from file & environment variables in Go with Viper

[viper](https://github.com/spf13/viper)<br>

<br><br>

### 2.2.1 why?

1. why file?
   development: easily specify default configuration for local development and testing

2. why env vars?
   development: easily override the default configurations when deploy with docker containers

3. why viper? <br>
   3.1 find, load, unmarshal config file<br>
   Json, TOML, YAML, ENV, INI<br><br>
   3.2 Read config from environment variables or flags <br>
   Override existing values, set default values<br><br>
   3.3 Read config from remote system<br>
   Etcd, consul<br><br>
   3.4 live watching and writing config file<br><br>
   Reread changed file, save any modifications<br><br>

<br><br>

### 2.2.2 viper setup

```bash
# (1) install viper
$ go get github.com/spf13/viper

# (2) create app.env file

# (3) create config.go

# (4) replace variables inside main.go

# (5) test in postman

# (6) override config before call make server
$ SERVER_ADDRESS=0.0.0.0:8081 make server
```

![imgs](./imgs/Xnip2023-01-05_12-27-23.jpg)

### 2.2.3 live watching

[live watching](https://github.com/spf13/viper)

### 2.2.4 reading from remote config system

[read from remote](https://github.com/spf13/viper)

<br><br>

## 2.3 Mock dB for testing HTTP API in Go and Achieve 100% coverage

<br><br>

### 2.3.1 why and how ?

#### 2.3.1.1 why mock?

1. Independent tests<br>
   Isolate tests data to avoid conflicts

2. Faster tests<br>
   Redue a lot of time talking to the database

3. 100% coverage<br>
   easily setup edge cases: unexpected errors

<br><br>

#### 2.3.1.2 how to mock?

1. use fake db: memory<br>
   implment a fake version of DB: store data in memory
   ![imgs](./imgs/Xnip2023-01-05_13-41-29.jpg)
2. use db stubs: gomock<br>
   Generate and build stubs that returns hard-coded values
   ![imgs](./imgs/Xnip2023-01-05_13-41-56.jpg)
   <br><br>

### 2.3.2 golang mock setup

```bash
# (1) install gomock
# https://github.com/golang/mock
# $ go get github.com/golang/mock/mockgen@v1.6.0 # for version < 1.16
$ go install github.com/golang/mock/mockgen@v1.6.0

# (2) validate mockgen
$ which mockgen

# (3) setup path in ~/.zshrc in most beginning
PATH = $PATH:~/go/bin

# (4) source
$ source ~/.zshrc

# (5) check again
$ which mockgen
```

![imgs](./imgs/Xnip2023-01-05_14-42-10.jpg)

```bash
# (1) add store interface

# (2) turn on emit_interface -> true in sqlc.yaml

# (3) regenerate
$ make sqlc

# (4) a new file named querier.go created

# (5) remove * from server.go, as it's no longer struct pointer type, but interface instead

# (6) as we have db.Store interface, now we can use it for mock
```

### 2.3.3 generate mock

```bash
# 1, specify path 2, specify interface name 3, sepcify the destination of generated output file
$ mockgen -destination db/mock/store.go github.com/RickZhou666/go-backend-service/db/sqlc Store

# 2. file generated to db/mock/store.go
type MockStore struct {}
type MockStoreMockRecorder struct {}

# 3. change the package name
$ mockgen -package mockdb -destination db/mock/store.go github.com/RickZhou666/go-backend-service/db/sqlc Store

# 4. add above cmd to MakeFile
```

### 2.3.4 write mock test

```bash
# (1) create account_test.go under api

# (2) create random account func

# (3) create new controller

# (4) create new mock store

# (5) setup store mock data

# (6) start test server and send request

# (7) execute unit test for test under api folder
$ go test -v --cover -count=1 ./api/...
```

![imgs](./imgs/Xnip2023-01-05_18-23-44.jpg)

2. `testing missing call`

```bash
# (1) comment store.GetAccount call in account.go

# (2) rerun test

```

![imgs](./imgs/Xnip2023-01-05_18-27-20.jpg)

```bash
# (3) as we define store.GetAccount func to be called 1 time. it's failed
```

![imgs](./imgs/Xnip2023-01-05_18-28-06.jpg)

3. `testing response boy`

```bash
# (1) define body matcher
```

![imgs](./imgs/Xnip2023-01-05_18-32-46.jpg)<br>

4. `test body not matched condition`<br>
   ![imgs](./imgs/Xnip2023-01-05_18-32-16.jpg)

<br><br>

### 2.3.5 Gin test mode setup

`click this run button`
![imgs](./imgs/Xnip2023-01-05_18-54-32.jpg)

<br><br>

1. gin.DebugMode

```go
gin.SetMode(gin.DebugMode)
```

![imgs](./imgs/Xnip2023-01-05_18-52-01.jpg)

<br><br>

2.  gin.TestMode

```go
gin.SetMode(gin.TestMode)
```

![imgs](./imgs/Xnip2023-01-05_18-53-41.jpg)

<br><br>

3.  gin.ReleaseMode

```go
gin.SetMode(gin.ReleaseMode)
```

![imgs](./imgs/Xnip2023-01-05_18-55-29.jpg)

<br><br>

4.  gin.EnvGinMode

```go
gin.SetMode(gin.EnvGinMode)
```

![imgs](./imgs/Xnip2023-01-05_18-55-50.jpg)

## 2.4 Implement transfer money API with a custom params validator

<br><br>

### 2.4.1 create transfer API

```bash
# (1) create struct

# (2) create func

# (3) add endpoint in server.go
```

<br><br>

### 2.4.2 customized validator for currency field

```bash
# (1) create currency validator

# (2) define currency util

# (3) setup currency check in server.go and replace in transfer.go handler
```

![imgs](./imgs/Xnip2023-01-05_22-37-18.jpg)

<br><br>

## 2.5 Add users table with unique & foreign key constaints in PostgreSQL

### 2.5.1 create new user table

```bash
# (1) create new users table
# Z - zero timezone or 00:00:00+00
password_changed_at timestamptz [not null, default: '0001-01-01 00:00:00Z']
```

![imgs](./imgs/Xnip2023-01-05_23-14-05.jpg)

```bash
# (2) one user can have mutilply currency account such USD account, RMB account
#     but cannot have multiple USD accounts, only one
```

<br><br>

### 2.5.2 migrate

<br><br>

### 2.5.2.1 replace init schema

1. replace init schema
2. reset db
3. run migrate up cmd<br>
   `this is not applicable in real world`

### 2.5.2.2 create new migration version

```bash
# (1) create new migration
$ migrate -help
$ migrate create -ext sql -dir db/migration -seq add_users
```

![imgs](./imgs/Xnip2023-01-05_23-35-59.jpg)

```bash
# (2) due to the dirty version, migrateup or migratedown will fail
```

![imgs](./imgs/Xnip2023-01-05_23-43-28.jpg)

```bash
# (3) overwrite dirty to FALSE

# (3) migratedown
$ make migratedown

# (4) migrate up
$ make migrateup

# (5) check users table
```

![imgs](./imgs/Xnip2023-01-05_23-53-07.jpg)<br>

![imgs](./imgs/Xnip2023-01-05_23-54-11.jpg)<br>

### 2.5.2.3 drop down

```bash
# (1) drop unique constraint for accounts table

# (2) drop foreign key constraint in similar way,
# check the name by clicking `info` button
```

![imgs](./imgs/Xnip2023-01-05_23-56-38.jpg)<br>

```bash
# (3) drop user table

# (4) setup makefile

# (4.1) migrate up to next 1 version
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

# (4.2) migrate down to next 1 version
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

# (5) after run migratedown1
#     users table disappear
#     constraint in accounts table disappear
#     schema_migrations version is 1
```

![imgs](./imgs/Xnip2023-01-06_00-03-07.jpg)

```bash
# (6) after run migrateup1
#     users table is here
#     constraint in accounts table is here
#     schema_migrations version is 2
```

<br><br>

<br><br>

## 2.6 How to handle DB errors in Golang correctly

<br><br>

## 2.7 How to securely store passwords? Hash password in Go with Bcrypt!

<br><br>

## 2.8 How to write stronger unit tests with a custom gomock matcher

<br><br>

## 2.9 Why PASETO is better than JWT for token-based authentication?

<br><br>

## 2.10 How to create and verify JWT & PASETO token in Golang

<br><br>

## 2.11 Implement login user API that returns PASETO or JWT access token in Go

<br><br>

## 2.12 Implement authenticaion middleware and authorization rules in Golang using Gin

<br><br><br>

# Section3: Deploying the application to production [Docker + Kubernetes + AWS]

<br><br><br>

# Section4: Advanced backend Topics [gRPC]

<br><br><br>

# Section5: Asynchronous processing with background workers [Asynq + Redis]
