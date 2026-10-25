.PHONY: build config up logs down ps

build:
	docker build --pull --tag weshnet-relay:0.5.0 .

config:
	docker compose config --quiet

up: config
	docker compose up -d

logs:
	docker compose logs --tail=100 -f

down ps:
	docker compose $@
