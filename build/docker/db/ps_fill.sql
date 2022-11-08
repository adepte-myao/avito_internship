INSERT INTO
    services (id, name)
SELECT
    series.series,
    concat('service_', series.series)
FROM
    generate_series(1, 50) as series;

INSERT INTO
    accounts (id, balance)
SELECT
    series.series,
    (series.series * 10) :: money
FROM
    generate_series(1, 5) as series;

INSERT INTO
    external_transfers_history (account_id, transfer_type, amount, record_time)
SELECT
    series.series,
    'deposit',
    (series.series * 10) :: money,
    now()
FROM
    generate_series(1, 5) as series;

INSERT INTO
    internal_transfers_history (sender_id, receiver_id, amount, record_time)
VALUES
    (1, 2, 10 :: money, now()),
    (2, 1, 10 :: money, now()),
    (1, 2, 10 :: money, now()),
    (2, 1, 10 :: money, now());

INSERT INTO
    reserves_history (
        account_id,
        service_id,
        order_id,
        total_cost,
        state,
        record_time,
        balance_after
    )
SELECT
    1,
    MOD(series.series, 50) + 1,
    1,
    0.01 :: money,
    'reserved',
    now(),
    (
        SELECT
            balance - (0.01 * series.series) :: money
        FROM
            accounts
        WHERE
            id = 1
    )
FROM
    generate_series(1, 10) as series;

INSERT INTO
    reserves_history (
        account_id,
        service_id,
        order_id,
        total_cost,
        state,
        record_time,
        balance_after
    )
SELECT
    1,
    MOD(series.series, 50) + 1,
    1,
    0.01 :: money,
    'accepted',
    now(),
    (
        SELECT
            balance - (0.01 * series.series) :: money
        FROM
            accounts
        WHERE
            id = 1
    )
FROM
    generate_series(1, 10) as series;

INSERT INTO
    reserves_history (
        account_id,
        service_id,
        order_id,
        total_cost,
        state,
        record_time,
        balance_after
    )
SELECT
    2,
    MOD(series.series, 50) + 6,
    1,
    0.02 :: money,
    'reserved',
    now(),
    (
        SELECT
            balance - (0.02 * series.series) :: money
        FROM
            accounts
        WHERE
            id = 2
    )
FROM
    generate_series(1, 10) as series;

INSERT INTO
    reserves_history (
        account_id,
        service_id,
        order_id,
        total_cost,
        state,
        record_time,
        balance_after
    )
SELECT
    2,
    MOD(series.series, 50) + 6,
    1,
    0.02 :: money,
    'cancelled',
    now(),
    (
        SELECT
            balance - 0.2 :: money + (0.02 * series.series) :: money
        FROM
            accounts
        WHERE
            id = 2
    )
FROM
    generate_series(1, 10) as series;

UPDATE
    accounts
SET
    balance = balance * 0.9
WHERE
    id = 1;