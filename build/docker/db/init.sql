CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    balance MONEY NOT NULL,
    CONSTRAINT non_negative_balance_check CHECK (balance >= 0 :: MONEY)
);

CREATE TABLE IF NOT EXISTS services (id SERIAL PRIMARY KEY, name TEXT);

CREATE TYPE reserve_state AS ENUM ('reserved', 'cancelled', 'accepted');

CREATE TABLE IF NOT EXISTS reserves_history (
    id BIGSERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts (id) NOT NULL,
    service_id INTEGER REFERENCES services (id) NOT NULL,
    order_id BIGINT NOT NULL,
    total_cost MONEY NOT NULL,
    state reserve_state NOT NULL,
    record_time TIMESTAMP WITH TIME ZONE NOT NULL,
    balance_after MONEY NOT NULL,
    CONSTRAINT unique_rule_reserves_history UNIQUE (account_id, service_id, order_id, state)
);

CREATE TABLE IF NOT EXISTS internal_transfers_history (
    id BIGSERIAL PRIMARY KEY,
    sender_id INTEGER REFERENCES accounts (id) NOT NULL,
    receiver_id INTEGER REFERENCES accounts (id) NOT NULL,
    amount MONEY NOT NULL,
    record_time TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TYPE transfer_type AS ENUM ('deposit', 'withdraw');

CREATE TABLE IF NOT EXISTS external_transfers_history (
    id BIGSERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts (id) NOT NULL,
    transfer_type transfer_type NOT NULL,
    amount MONEY NOT NULL,
    record_time TIMESTAMP WITH TIME ZONE NOT NULL
);
