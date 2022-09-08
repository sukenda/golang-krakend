package client

import (
	"fmt"
	"github.com/sukenda/golang-krakend/grpc-proto/proto"

	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	Client proto.AuthServiceClient
}

func InitAuthServiceClient(url string) AuthServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := AuthServiceClient{
		Client: proto.NewAuthServiceClient(cc),
	}

	return c
}
