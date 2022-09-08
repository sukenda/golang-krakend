package services

import (
	"context"
	database "github.com/sukenda/golang-krakend/product-service/database"
	"github.com/sukenda/golang-krakend/product-service/model"
	"github.com/sukenda/golang-krakend/proto"
	"net/http"
)

type ProductService struct {
	proto.UnimplementedProductServiceServer
	Postgre database.Postgre
}

func (s *ProductService) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	var product model.Product

	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price

	if result := s.Postgre.GormBD.Create(&product); result.Error != nil {
		return &proto.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}

	return &proto.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}, nil
}

func (s *ProductService) FindOne(ctx context.Context, req *proto.FindOneRequest) (*proto.FindOneResponse, error) {
	var product model.Product

	if result := s.Postgre.GormBD.First(&product, req.Id); result.Error != nil {
		return &proto.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	data := &proto.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &proto.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (s *ProductService) DecreaseStock(ctx context.Context, req *proto.DecreaseStockRequest) (*proto.DecreaseStockResponse, error) {
	var product model.Product

	if result := s.Postgre.GormBD.First(&product, req.Id); result.Error != nil {
		return &proto.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	if product.Stock <= 0 {
		return &proto.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock too low",
		}, nil
	}

	var log model.StockDecreaseLog

	if result := s.Postgre.GormBD.Where(&model.StockDecreaseLog{OrderId: req.OrderId}).First(&log); result.Error == nil {
		return &proto.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock already decreased",
		}, nil
	}

	product.Stock = product.Stock - 1

	s.Postgre.GormBD.Save(&product)

	log.OrderId = req.OrderId
	log.ProductRefer = product.Id

	s.Postgre.GormBD.Create(&log)

	return &proto.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}
