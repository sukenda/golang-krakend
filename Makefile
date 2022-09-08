compose-build:
	docker-compose -f docker-compose.yaml build

compose-up:
	docker-compose -f docker-compose.yaml build
	docker-compose -f docker-compose.yaml up

compose-down:
	docker-compose -f docker-compose.yaml down

compose-restart:
	docker-compose -f docker-compose.yaml down
	docker-compose -f docker-compose.yaml up