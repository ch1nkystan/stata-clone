create table if not exists addresses
(
    id          serial
        constraint address_pk primary key,

    blockchain  varchar(32) default ''    not null,
    address_key text        default ''    not null unique,
    address     text        default ''    not null,

    bid         varchar(32) default ''    not null,

    created_at  timestamp   default now() not null,
    updated_at  timestamp   default now() not null
);