create table if not exists users(
    id uuid primary key default gen_random_uuid()
);

insert into users(id) values('00000000-0000-0000-0000-000000000001');