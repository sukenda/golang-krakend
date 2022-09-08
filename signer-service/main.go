package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/sukenda/golang-krakend/grpc-proto/proto"
	"github.com/sukenda/golang-krakend/signer-service/config"
	"github.com/sukenda/golang-krakend/signer-service/services"
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

	signerService := services.SignerService{}

	// Create your protocol servers.
	grpcServer := grpc.NewServer()
	proto.RegisterSignerServiceServer(grpcServer, &signerService)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = proto.RegisterSignerServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%v", c.Port), opts)
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
