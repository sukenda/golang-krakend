gomod:
	rm -rf vendor && rm -rf go.sum && go mod tidy && go mod vendor

run:
	go run main.go

build:
	docker build -t signer-service .

