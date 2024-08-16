create table if not exists transactions
(
    id         serial
        constraint transactions_pk primary key,
    user_id    int         default 0     not null,

    blockchain varchar(32) default ''    not null,
    tx_hash    text        default ''    not null,
    tx_key     text        default ''    not null unique,

    amount     float4      default 0     not null,
    price      float4      default 0     not null,

    created_at timestamp   default now() not null,
    updated_at timestamp   default now() not null
);