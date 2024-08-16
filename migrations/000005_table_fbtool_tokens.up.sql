create table if not exists fbtool_tokens
(
    id            serial
        constraint fbtool_tokens_pk primary key,
    token         varchar(255) default ''    not null unique,
    active        boolean      default true  not null,
    days_to_fetch integer      default 10    not null,

    created_at    timestamp    default now() not null,
    fetched_at    timestamp    default now() not null
);