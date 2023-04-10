create table accounts (
    id bigserial not null primary key,
    creation_date date not null default CURRENT_DATE,
    account_type varchar not null,
    balance numeric(11, 3) default 0
);