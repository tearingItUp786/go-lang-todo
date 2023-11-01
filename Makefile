docker\:up:
	docker compose up
docker\:down:
	docker compose down
dev:
	air kill & npm run build:watch
