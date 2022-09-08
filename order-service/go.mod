module github.com/sukenda/golang-krakend/order-service

go 1.16

require (
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3
	github.com/sirupsen/logrus v1.9.0
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/viper v1.12.0
	github.com/sukenda/golang-krakend/grpc-proto v0.0.0-20220908014605-a33e5a9b45ad
	google.golang.org/grpc v1.49.0
	gorm.io/driver/postgres v1.3.8
	gorm.io/gorm v1.23.8
)