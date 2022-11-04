CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    balance MONEY NOT NULL,
    CONSTRAINT nonNegativeBalanceCheck CHECK (balance >= 0::MONEY)
);
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY, 
    name TEXT
);
CREATE TYPE reserve_state AS ENUM ('reserved', 'cancelled', 'accepted');
CREATE TABLE IF NOT EXISTS reserves_history (
    id BIGSERIAL PRIMARY KEY,
    accountID INTEGER REFERENCES accounts (id) NOT NULL,
    serviceID INTEGER REFERENCES services (id) NOT NULL,
    orderID BIGINT NOT NULL,
    totalCost MONEY NOT NULL,
    state reserve_state NOT NULL,
    record_time TIMESTAMP WITH TIME ZONE NOT NULL,
    balanceAfter MONEY NOT NULL
);
CREATE TYPE transfer_type AS ENUM ('deposit', 'withdraw');
CREATE TABLE IF NOT EXISTS custom_transfers_history (
    id BIGSERIAL PRIMARY KEY,
    accountID INTEGER REFERENCES accounts (id) NOT NULL,
    otherAccountID INTEGER REFERENCES accounts (id) DEFAULT NULL,
    transferType transfer_type NOT NULL,
    amount MONEY NOT NULL,
    record_time TIMESTAMP WITH TIME ZONE NOT NULL,
    balanceAfter MONEY NOT NULL
);