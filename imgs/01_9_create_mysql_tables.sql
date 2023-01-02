-- (1) create table
CREATE TABLE IF NOT EXISTS accounts(  
    id INT NOT NULL AUTO_INCREMENT,
    owner VARCHAR(255) NOT NULL,
    balance INT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
) ENGINE=InnoDB;

-- (2) drop table
drop table accounts;

-- (3) insert records
INSERT INTO accounts(owner,balance,currency) VALUES ('rick', 100, 'USD');
INSERT INTO accounts(owner,balance,currency) VALUES ('eve', 100, 'USD');

-- (4) select all records
select * from accounts;

-- (5) check table attributes
DESC accounts;


-- (4) show database;
SHOW databases;
-- (5) create database;
create database simple_bank;

-- (6) switch databases;
USE simple_bank;


CREATE TABLE customers 
(a INT, b CHAR (20), INDEX (a)) 
ENGINE=InnoDB;


CREATE TABLE movies(title VARCHAR(50) NOT NULL,
genre VARCHAR(30) NOT NULL,
director VARCHAR(60) NOT NULL,
release_year INT NOT NULL,
PRIMARY KEY(title));


CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);