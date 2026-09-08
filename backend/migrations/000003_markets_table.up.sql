create table if not exists markets(
    id uuid primary key default gen_random_uuid(),
    title text not null,
    api_link text not null
);

insert into markets(id, title, api_link) values('00000000-0000-0000-0000-000000000042', 'MegaMarket', 'market:8081');