package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/sukenda/golang-krakend/auth-service/config"
	db "github.com/sukenda/golang-krakend/auth-service/database"
	"github.com/sukenda/golang-krakend/auth-service/services"
	"github.com/sukenda/golang-krakend/auth-service/utils"
	"github.com/sukenda/golang-krakend/grpc-proto/proto"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)
	jwt := utils.JwtWrapper{
		Kid:             "bluebird.id",
		Issuer:          "bluebird.id",
		SecretKey:       c.JWTSecretKey,
		ExpirationHours: 24,
	}

	// Create the main listener.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", c.Port))
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux.
	cMux := cmux.New(listener)

	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcL := cMux.Match(cmux.HTTP2(), cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc"))
	httpL := cMux.Match(cmux.HTTP1Fast())

	//prvKey, err := ioutil.ReadFile("./config/id_rsa")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//pubKey, err := ioutil.ReadFile("./config/id_rsa.pub")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	jwtToken := utils.NewJWT(nil, nil)
	authService := services.AuthService{
		Database:   h,
		JwtWrapper: jwt,
		JWT:        jwtToken,
	}

	// Create your protocol servers.
	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, &authService)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = proto.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%v", c.Port), opts)
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Addr:    c.Port,
		Handler: gwmux,
	}

	// Use the muxed listeners for your servers.
	go grpcServer.Serve(grpcL)
	go httpServer.Serve(httpL)

	// Start serving!
	err = cMux.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
