package service

import (
	"context"
	"github.com/sukenda/golang-krakend/grpc-proto/client"
	"github.com/sukenda/golang-krakend/grpc-proto/proto"
	database "github.com/sukenda/golang-krakend/order-service/database"
	models "github.com/sukenda/golang-krakend/order-service/model"
	"net/http"
)

type OrderService struct {
	proto.UnimplementedOrderServiceServer
	Database       database.Handler
	ProductService client.ProductServiceClient
}

func (s *OrderService) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	product, err := s.ProductService.FindOne(req.ProductId)

	if err != nil {
		return &proto.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if product.Status >= http.StatusNotFound {
		return &proto.CreateOrderResponse{Status: product.Status, Error: product.Error}, nil
	} else if product.Data.Stock < req.Quantity {
		return &proto.CreateOrderResponse{Status: http.StatusConflict, Error: "Stock too less"}, nil
	}

	order := models.Order{
		Price:     product.Data.Price,
		ProductId: product.Data.Id,
		UserId:    req.UserId,
	}

	s.Database.DB.Create(&order)

	res, err := s.ProductService.DecreaseStock(req.ProductId, order.Id)

	if err != nil {
		return &proto.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if res.Status == http.StatusConflict {
		s.Database.DB.Delete(&models.Order{}, order.Id)

		return &proto.CreateOrderResponse{Status: http.StatusConflict, Error: res.Error}, nil
	}

	return &proto.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}
