create table if not exists products(
    id uuid primary key default gen_random_uuid(),
    market_id uuid not null references markets(id) on delete cascade,
    title text not null,
    price int not null
);

create index if not exists ind_product on products(id);