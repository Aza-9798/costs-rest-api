alter table transactions
add constraint check_transactions_source_destination
    check (transactions.source != transactions.destination);