create table if not exists carts(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    in_process boolean default false,
    ordered_at timestamp default null
);