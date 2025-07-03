create table if not exists pixel_links
(
    id               serial
        constraint pixel_links_pk primary key,

    fb_access_marker varchar(255) default ''    not null,
    fb_pixel_id      bigint       default 0     not null,
    fbc              varchar(255) default ''    not null,
    fbp              varchar(255) default ''    not null,

    deeplink_id      int          default 0     not null,

    invite_uuid      uuid                       not null default '00000000-0000-0000-0000-000000000000',

    created_at       timestamp    default now() not null
);

alter table users add column invite_uuid uuid not null default '00000000-0000-0000-0000-000000000000';

CREATE INDEX IF NOT EXISTS idx_pixel_links_deeplink_id ON pixel_links(deeplink_id);
CREATE INDEX IF NOT EXISTS idx_pixel_links_fb_access_marker ON pixel_links(fb_access_marker);
CREATE INDEX IF NOT EXISTS idx_pixel_links_fb_pixel_id ON pixel_links(fb_pixel_id);