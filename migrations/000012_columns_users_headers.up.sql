alter table users
    add column if not exists ip           varchar(64)  default '' not null,
    add column if not exists user_agent   varchar(255) default '' not null,
    add column if not exists country_code varchar(2)   default '' not null,
    add column if not exists os_name      varchar(255) default '' not null,
    add column if not exists device_type  varchar(255) default '' not null;