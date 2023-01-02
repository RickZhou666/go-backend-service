-- 1. replicate deadlock
-- Tx1: transfer $10 from account 1 to account 2
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;

ROLLBACK;

-- Tx2: transfer $10 from account 2 to account 1
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;
UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;

ROLLBACK;


-- check lock
SELECT a.application_name,
         l.relation::regclass,
         l.transactionid,
         l.mode,
         l.locktype,
         l.GRANTED,
         a.usename,
         a.query,
         a.query_start,
         a.pid
FROM pg_stat_activity a
JOIN pg_locks l ON l.pid = a.pid
WHERE a.application_name = 'psql'
ORDER BY a.query_start;

-- 2. avoid deadlock

-- Tx1: transfer $10 from account 1 to account 2
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;

COMMIT;

-- Tx2: transfer $10 from account 2 to account 1
BEGIN;

UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;

COMMIT;


-- 1. get at most 5 records from accounts
-- offset = 10, skip 10 records
-- limits = 5, as most 5 records returned
SELECT * FROM accounts
ORDER BY id
LIMIT 5
OFFSET 10;

-- 2. entries
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = 1
ORDER BY id
LIMIT 5
OFFSET 10;


-- 3. select all account_id = 1
SELECT * FROM entries e
WHERE e.account_id = 1
ORDER BY ID;

-- 3. return 17 records
SELECT COUNT(*) FROM entries e
WHERE e.account_id = 1;

-- 4. [FALSE] return 0. as result only has 1 row
SELECT COUNT(*) FROM entries e
WHERE e.account_id = 1
LIMIT 5
OFFSET 10;

-- 5. add id constraints
SELECT COUNT(*) FROM entries e
WHERE e.account_id = 1
    AND e.id >= 185;