create table if not exists deeplinks
(
    id                   serial
        constraint deeplinks_pk primary key,

    bot_id               int          default 0     not null,

    referral_telegram_id bigint                     not null,
    hash                 varchar(255) default ''    not null,
    label                varchar(255) default ''    not null,

    active               boolean      default true  not null,

    created_at           timestamp    default now() not null,
    updated_at           timestamp    default now() not null
);

alter table deeplinks
    add constraint bot_deeplink_hash_unique UNIQUE (bot_id, hash);