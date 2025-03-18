
-- Таблицы
CREATE TABLE IF NOT EXISTS 
    accounts (
        number          TEXT NOT NULL PRIMARY KEY,
        full_name       TEXT NOT NULL,
        exparation_date TEXT NOT NULL,
        cvv             INTEGER NOT NULL,
        balance         BIGINT NOT NULL CHECK(balance >= 0) DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS 
    transfers (
        id          INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id   TEXT NOT NULL REFERENCES accounts(account_id),
        receiver_id TEXT NOT NULL CHECK(receiver_id != sender_id) REFERENCES accounts(account_id),
        amount      BIGINT NOT NULL CHECK(amount > 0),
        status      TEXT NOT NULL
    );
