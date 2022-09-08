protogen:
	@echo "Loading ..."
	@go get google.golang.org/grpc@v1.48.0
	@go get google.golang.org/protobuf@v1.28.1
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.11.2
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.3

	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/order.proto
	@protoc --grpc-gateway_out=logtostderr=true,grpc_api_configuration=proto/order.yaml:proto/. proto/order.proto
	@protoc --openapiv2_out=logtostderr=true,grpc_api_configuration=proto/order.yaml:. proto/order.proto

	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/product.proto
	@protoc --grpc-gateway_out=logtostderr=true,grpc_api_configuration=proto/product.yaml:proto/. proto/product.proto
	@protoc --openapiv2_out=logtostderr=true,grpc_api_configuration=proto/product.yaml:. proto/product.proto

	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/auth.proto
	@protoc --grpc-gateway_out=logtostderr=true,grpc_api_configuration=proto/auth.yaml:proto/. proto/auth.proto
	@protoc --openapiv2_out=logtostderr=true,grpc_api_configuration=proto/auth.yaml:. proto/auth.proto
	@echo "Done"

compose-build:
	docker-compose -f docker-compose.yaml build

compose-up:
	docker-compose -f docker-compose.yaml build
	docker-compose -f docker-compose.yaml up

compose-down:
	docker-compose -f docker-compose.yaml down

compose-restart:
	docker-compose -f docker-compose.yaml down
	docker-compose -f docker-compose.yaml up -d