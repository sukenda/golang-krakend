module github.com/sukenda/golang-krakend/auth-service

go 1.16

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3
	github.com/sirupsen/logrus v1.6.0
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/viper v1.12.0
	github.com/sukenda/golang-krakend/grpc-proto v0.0.0-20220908013721-d1d01bc5ca21
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/postgres v1.3.8
	gorm.io/gorm v1.23.8
)
