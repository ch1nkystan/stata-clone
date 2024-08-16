create table if not exists bots
(
    id           serial
        constraint bots_pk primary key,
    api_key      varchar(255) default ''    not null unique,
    bot_token    varchar(255) default ''    not null unique,
    bot_username varchar(255) default ''    not null,
    bot_type     varchar(255) default ''    not null,

    active       boolean      default true  not null,
    bid          varchar(32)  default ''    not null,

    created_at   timestamp    default now() not null,
    updated_at   timestamp    default now() not null
);