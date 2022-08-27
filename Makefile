compose-build:
	docker-compose -f docker-compose.yaml build

compose-up:
	docker-compose -f docker-compose.yaml up -d

compose-down:
	docker-compose -f docker-compose.yaml down