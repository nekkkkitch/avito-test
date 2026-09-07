-- name: GetMarkets :many
select * from markets;

-- name: GetProductCard :one
select * from products where id = $1;