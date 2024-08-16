create table if not exists fbtool_campaigns_stats
(
    id                serial
        constraint fbtool_campaigns_stats_pk primary key,

    fbtool_account_id int          default 0     not null,

    campaign_name     varchar(255) default ''    not null,
    campaign_id       varchar(255) default ''    not null,

    status            varchar(255) default ''    not null,
    effective_status  varchar(255) default ''    not null,

    impressions       int          default 0     not null,
    clicks            int          default 0     not null,
    spend             float        default 0     not null,
    date              timestamp    default now() not null,

    created_at        timestamp    default now() not null
);

alter table fbtool_campaigns_stats
    add constraint daily_campaign_stat UNIQUE (fbtool_account_id, campaign_id, date);