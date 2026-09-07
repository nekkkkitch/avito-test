start:
	sudo docker compose -f docker-compose.yml --env-file config.env up -d

rebuild:
	sudo docker rmi deploy-backend -f && make start
