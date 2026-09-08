-- name: GetCart :one
select * from carts where user_id = $1 and ordered_at is null;

-- name: GetCartProducts :many
select * from products p join carts_content c on p.id = c.product_id where c.cart_id = $1;

-- name: GetProductsForOrder :many
select * from products p join carts_content c on p.id = c.product_id join markets m on p.market_id = m.id where c.cart_id = $1; 

-- name: CreateCart :one
insert into carts(user_id) values($1) returning *;

-- name: AddToCart :exec
insert into carts_content(cart_id, product_id) values($1, $2);

-- name: DeleteFromCart :exec
delete from carts_content 
where id in (select id from carts_content c where c.cart_id = $1 and c.product_id = $2 limit 1);

-- name: OrderStartProcess :exec
update carts set in_process = true where id = $1 and in_process = false;

-- name: FinishOrder :exec
update carts set ordered_at = now() where id = $1;