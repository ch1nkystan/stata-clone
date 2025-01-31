create table if not exists users
(
    id                       serial
        constraint users_pk primary key,
    bot_id                   int          default 0       not null,
    deeplink_id              int          default 0       not null,

    depot_channel_hash       varchar(16)  default ''      not null,
    telegram_channel_id      bigint                       not null,
    telegram_channel_url     varchar(255) default ''      not null,

    telegram_id              bigint       default 0       not null,
    first_name               varchar(255) default ''      not null,
    last_name                varchar(255) default ''      not null,
    username                 varchar(255) default ''      not null,
    forward_sender_name      varchar(255) default ''      not null,

    ip                       varchar(39)  default ''      not null,
    user_agent               varchar(255) default ''      not null,
    country_code             varchar(2)   default ''      not null,
    os_name                  varchar(255) default ''      not null,
    device_type              varchar(255) default ''      not null,

    is_bot                   boolean      default false   not null,
    is_premium               boolean      default false   not null,
    language_code            varchar(8)   default ''      not null,

    seen                     int          default 0       not null,
    active                   boolean      default true    not null,
    event_created            varchar(255) default ''      not null,

    deposits_total           int          default 0       not null,
    deposits_sum             float4       default 0       not null,
    deposited                boolean      default false   not null,
    deposited_at             timestamp    default now()   not null,

    mailing_state            varchar(32)  default 'ready' not null,
    mailing_state_updated_at timestamp    default now()   not null,
    mailing_failed_attempts  int          default 0       not null,

    messaged_at              timestamp    default now()   not null,
    created_at               timestamp    default now()   not null,
    updated_at               timestamp    default now()   not null
);