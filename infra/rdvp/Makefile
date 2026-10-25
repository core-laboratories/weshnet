.PHONY: all config up logs down ps genkey ip

all: up ps logs

config:
	@test -f .env || { echo "Copy .env.example to .env and configure it first." >&2; exit 1; }
	@set -a; . ./.env; set +a; \
		for name in RDVP_IMAGE RDVP_PK ANNOUNCE_SERVER EMITTER_SECRET_KEY EMITTER_LICENSE; do \
			eval "value=\$$$${name}"; \
			test -n "$$value" || { echo "$$name is required in .env" >&2; exit 1; }; \
		done
	docker compose config --quiet

up: config
	docker compose up -d

logs:
	docker compose logs --tail=100 -f

down ps:
	docker compose $@

genkey:
	@test -f .env || cp .env.example .env
	@set -a; . ./.env; set +a; \
		test -n "$$RDVP_IMAGE" || { echo "RDVP_IMAGE is required in .env" >&2; exit 1; }; \
		key="$$(docker run --rm --entrypoint rdvp "$$RDVP_IMAGE" genkey)"; \
		test -n "$$key"; \
		if grep -q '^RDVP_PK=' .env; then \
			sed -i.bak "s|^RDVP_PK=.*|RDVP_PK=$$key|" .env && rm -f .env.bak; \
		else \
			printf '\nRDVP_PK=%s\n' "$$key" >> .env; \
		fi

ip:
	curl --fail --silent --show-error https://ifconfig.co
