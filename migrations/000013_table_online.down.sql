drop table if exists online;

drop index idx_users_messaged_at ON users(messaged_at);