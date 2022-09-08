package client

import (
	"fmt"
	"github.com/sukenda/golang-krakend/grpc-proto/proto"

	"google.golang.org/grpc"
)

type OrderServiceClient struct {
	Client proto.OrderServiceClient
}

func InitOrderServiceClient(url string) OrderServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := OrderServiceClient{
		Client: proto.NewOrderServiceClient(cc),
	}

	return c
}
