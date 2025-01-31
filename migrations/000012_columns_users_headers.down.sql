alter table users
    drop column if exists ip,
    drop column if exists user_agent,
    drop column if exists country_code,
    drop column if exists os_name,
    drop column if exists device_type;