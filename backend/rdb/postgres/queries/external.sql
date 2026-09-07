-- name: SetProduct :exec
insert into products(id, market_id, title, price) values($1, $2, $3, $4)
on conflict (id) do 
update set market_id = $2, title = $3, price = $4;

-- name: DeleteProduct :exec
delete from products where id = $1;