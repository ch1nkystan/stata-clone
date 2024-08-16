create table if not exists prices
(
    id         serial
        constraint price_id primary key,
    ticker     varchar(32) default ''    not null,
    price      float4      default 0     not null,
    created_at timestamp   default now() not null
);