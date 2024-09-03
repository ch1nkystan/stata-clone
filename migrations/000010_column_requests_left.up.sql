alter table fbtool_tokens
    add column if not exists requests_left int not null default 0;
