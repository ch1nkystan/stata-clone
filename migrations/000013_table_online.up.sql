create table if not exists snapshots
(
    id            serial
        constraint snapshots_pk primary key,
    snapshot varchar(32) default '' not null,
    bot_id        integer                not null,
    users         integer                not null,
    created_at    timestamp              not null
);

create index idx_users_messaged_at ON users(messaged_at);

alter table snapshots
    add constraint snapshot_cat_unique UNIQUE (bot_id, snapshot, created_at);