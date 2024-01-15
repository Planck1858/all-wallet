-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.users
(
    id          integer generated always as identity,
    telegram_id integer unique not null,
    total_money numeric(20,10) not null,
    currency    varchar(10) not null,
    created_at  timestamp without time zone not null,
    PRIMARY KEY(id)
);

CREATE TABLE public.account
(
    id       integer generated always as identity,
    name     varchar(255) not null,
    type     varchar(255) not null,
    balance  numeric(20,10) not null,
    currency varchar(10) not null,
    user_id  integer not null,
    PRIMARY KEY(id),
    CONSTRAINT fk_account_to_users
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);

CREATE TABLE public.spending
(
    id       integer generated always as identity,
    total    numeric(20,10) not null,
    currency varchar(10) not null,
    user_id  integer not null,
    PRIMARY KEY(id),
    CONSTRAINT fk_spending_to_users
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);

CREATE TABLE public.spending_record
(
    id          integer generated always as identity,
    amount      numeric(20,10) not null,
    -- TODO: tag         varchar(255), /* array? */
    currency    varchar(10) not null,
    spending_id integer not null,
    created_at  timestamp without time zone not null,
    CONSTRAINT fk_spending_record_to_spending
        FOREIGN KEY(spending_id)
            REFERENCES spending(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.spending_record;
DROP TABLE public.spending;
DROP TABLE public.account;
DROP TABLE public.users;
-- +goose StatementEnd
