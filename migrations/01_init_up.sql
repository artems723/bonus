CREATE TABLE IF NOT EXISTS "user"
(
    id            bigserial primary key,
    login         text   not null,
    password_hash text   not null,
    CONSTRAINT    login_unique UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS "order"
(
    id          bigserial primary key,
    user_id     bigint         not null,
    number      varchar(255)   not null,
    status      varchar(10)    not null,
    accrual     numeric(12, 2) not null default 0,
    uploaded_at Timestamp      not null,
    CONSTRAINT  number_unique UNIQUE (number)
);

CREATE TABLE IF NOT EXISTS "balance"
(
    id          bigserial primary key,
    user_id     bigint         not null,
    ordernum    varchar(255)   not null,
    debit       numeric(12, 2) not null,
    credit      numeric(12, 2) not null,
    created_at  Timestamp      not null,
    CONSTRAINT  ordernum_unique UNIQUE (ordernum)
);