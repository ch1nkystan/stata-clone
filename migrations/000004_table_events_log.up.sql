create table if not exists events_log
(
    id                   serial
        constraint deposit_log_pk primary key,
    event_type           varchar(255) default ''    not null,

    user_id              int          default 0     not null,
    reporter_telegram_id bigint                     not null,
    telegram_id          bigint       default 0     not null,

    created_at           timestamp    default now() not null,
    updated_at           timestamp    default now() not null
);