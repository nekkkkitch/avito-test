start:
	sudo docker compose -f docker-compose.yml --env-file config.env up -d

rebuild:
	sudo docker rmi deploy-backend -f && make start

swag_main:
	swag init --parseDependencyLevel 1 -d ../backend/server --generalInfo ../server/server.go -o ../backend/api/
