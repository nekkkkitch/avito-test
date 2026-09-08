start:
	sudo docker compose -f docker-compose.yml --env-file config.env up -d

rebuild_main:
	sudo docker rmi deploy-backend -f && make start

rebuild_market:
	sudo docker rmi deploy-market -f && make start

full_rebuild:
	sudo docker rmi deploy-market -f && sudo docker rmi deploy-backend -f && sudo docker container rm -f deploy-postgres-1 && make start

swag_main:
	swag init --parseDependencyLevel 1 -d ../backend/server --generalInfo ../server/server.go -o ../backend/api/

swag_market:
	swag init --parseDependencyLevel 1 -d ../market/server --generalInfo ../server/server.go -o ../market/api/