CREATE TABLE IF NOT EXISTS users
(
    id            bigserial primary key,
    login         varchar(255)   not null,
    password_hash varchar(255)   not null,
    CONSTRAINT    login_unique UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS orders
(
    id          bigserial primary key,
    user_id     bigint         not null,
    number      varchar(255)   not null,
    status      varchar(10)    not null,
    accrual     numeric(12, 2),
    uploaded_at Timestamp      not null,
    CONSTRAINT  number_unique UNIQUE (number)
);

CREATE TABLE IF NOT EXISTS balances
(
    id          bigserial primary key,
    user_id     bigint         not null,
    ordernum    varchar(255)   not null,
    debit       numeric(12, 2),
    credit      numeric(12, 2),
    created_at  Timestamp      not null,
    CONSTRAINT  ordernum_unique UNIQUE (ordernum)
);