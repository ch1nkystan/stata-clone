alter table bots add column trace_uuid varchar(36) default '00000000-0000-0000-0000-000000000000' not null;
alter table bots add column binding boolean default false not null;