create table if not exists carts_content(
    id serial,
    cart_id uuid not null references carts(id) on delete cascade,
    product_id uuid not null references products(id) on delete cascade
);

create index if not exists ind_cart_content on carts_content(cart_id, product_id);