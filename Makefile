docker\:up:
	docker compose up

docker\:down:
	docker compose down

dev:
	air kill & npm run build:watch

db-proxy\:start:
	cloudflared access tcp --hostname pg.taranveerbains.ca --url localhost:6432 &

db-proxy\:stop:
	pkill cloudflared
