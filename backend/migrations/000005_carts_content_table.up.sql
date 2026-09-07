create table if not exists carts_content(
    cart_id uuid not null references carts(id) on delete cascade,
    product_id uuid not null references products(id) on delete cascade
);