alter table users
    add column if not exists subscribed      boolean   default false not null,
    add column if not exists subscribed_at   timestamp default now() not null,
    add column if not exists unsubscribed_at timestamp default now() not null;

create index idx_users_subscribed_at ON users(subscribed_at);
create index idx_users_unsubscribed_at ON users(unsubscribed_at);