create table transactions (
    id bigserial not null primary key,
    creation_date date not null default CURRENT_DATE,
    transaction_date date not null default CURRENT_DATE,
    description varchar(1000),
    source bigint not null references accounts(id),
    destination bigint not null references accounts(id),
    amount numeric(11, 3) not null
);