create table if not exists fbtool_accounts
(
    id                  serial
        constraint fbtool_accounts_pk primary key,

    token_id            int          default 0     not null,
    fbtool_account_id   int          default 0     not null unique,
    fbtool_account_name varchar(255) default ''    not null,

    active              boolean      default true  not null,

    created_at          timestamp    default now() not null,
    fetched_at          timestamp    default now() not null
);