drop index idx_users_subscribed_at;
drop index idx_users_unsubscribed_at;

alter table users
    drop column if exists subscribed,
    drop column if exists subscribed_at,
    drop column if exists unsubscribed_at;