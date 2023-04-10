alter table accounts
add column user_id int not null;

alter table accounts
add foreign key(user_id) references users(id);