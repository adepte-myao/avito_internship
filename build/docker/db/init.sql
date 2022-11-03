CREATE TABLE IF NOT EXISTS user_balances (
    id SERIAL PRIMARY KEY,
    balance MONEY NOT NULL,
    CONSTRAINT nonNegativeBalanceCheck 
        CHECK ( balance >= 0 )
);

CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS active_reserves (
    id SERIAL PRIMARY KEY,
    balanceID INTEGER REFERENCES user_balances (id) NOT NULL,
    serviceID INTEGER REFERENCES services (id) NOT NULL,
    orderID BIGINT NOT NULL,
    totalCost MONEY NOT NULL
);

CREATE TYPE reserve_state AS ENUM ('reserved', 'cancelled', 'accepted');

CREATE TABLE IF NOT EXISTS reserves_history (
    id BIGSERIAL PRIMARY KEY,
    balanceID INTEGER REFERENCES user_balances (id) NOT NULL,
    serviceID INTEGER REFERENCES services (id) NOT NULL,
    orderID BIGINT NOT NULL,
    totalCost MONEY NOT NULL,
    state reserve_state NOT NULL,
    record_time TIMESTAMP WITH TIME ZONE NOT NULL,
    balanceAfter MONEY NOT NULL
);

CREATE TYPE transfer_type AS ENUM ('crediting', 'debiting');

CREATE TABLE IF NOT EXISTS custom_transfers_history (
    id BIGSERIAL PRIMARY KEY,
    balanceID INTEGER REFERENCES user_balances (id) NOT NULL,
    transferType transfer_type NOT NULL,
    amount MONEY NOT NULL,
    record_time TIMESTAMP WITH TIME ZONE NOT NULL,
    balanceAfter MONEY NOT NULL
);