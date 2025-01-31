create table if not exists online_history (
    bot_id              integer     not null,
    interval_start      timestamp   not null,
    interval_end        timestamp   not null, 
    active_users_count  integer     not null
);

create index idx_online_history_interval on online_history (bot_id, interval_start, interval_end);