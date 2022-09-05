package services

import (
	"context"
	database "github.com/sukenda/golang-krakend/auth-service/database"
	models "github.com/sukenda/golang-krakend/auth-service/model"
	"github.com/sukenda/golang-krakend/auth-service/proto"
	"github.com/sukenda/golang-krakend/auth-service/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type AuthService struct {
	proto.UnimplementedAuthServiceServer
	Database   database.Handler
	JwtWrapper utils.JwtWrapper
}

func (s *AuthService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	var user models.User

	if result := s.Database.GormDb.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return &proto.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)

	s.Database.GormDb.Create(&user)

	return &proto.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user models.User

	if result := s.Database.GormDb.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return &proto.LoginResponse{
			Status: http.StatusNotFound,
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		return &proto.LoginResponse{
			Status: http.StatusForbidden,
		}, nil
	}

	token, refresh, exp, err := s.JwtWrapper.GenerateToken(user)
	if err != nil {
		return &proto.LoginResponse{
			Status: http.StatusBadGateway,
		}, nil
	}

	return &proto.LoginResponse{
		Status:       http.StatusOK,
		AccessToken:  token,
		RefreshToken: refresh,
		Exp:          exp,
	}, nil
}

func (s *AuthService) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, err := s.JwtWrapper.ValidateToken(req.Token)

	if err != nil {
		return &proto.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User

	if result := s.Database.GormDb.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		return &proto.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &proto.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}

func (s *AuthService) JWKValidate(ctx context.Context, req *emptypb.Empty) (*proto.JWKValidateResponse, error) {
	return &proto.JWKValidateResponse{
		Keys: []*proto.JWK{
			{
				Alg: "HS256",
				Kid: "bluebird.id",
				Use: "enc",
				Kty: "EC",
				Crv: "P-256",
				X:   "2xlQ_IYvfS1cUXTb3orNSxIJ8B8b7EcYDgdMkpyLmp8",
				Y:   "MDSo0odLkrtNN13Hb7pAv1pOnGqYmqb0FgJpRTUUW0s",
			},
			{
				Alg: "HS256",
				Kid: "bluebird.id",
				Kty: "EC",
				Crv: "P-256",
				X:   "nqqhD6gwxnGtj1U3cx_em4qEaI2EUXTnWP4W8d5VpOI",
				Y:   "EzgfVjlaSluvcfWpmqYvGO_RObDPXCCIWmg7nIiFGew",
			},
			{
				Alg: "HS256",
				Kid: "bluebird.id",
				Kty: "EC",
				Crv: "P-256",
				X:   "YRdzHwnDsxbUqhUYH-Hrz8R8vVXc_slU3P-k6uGAvyQ",
				Y:   "inzssI5Sjk1MyXKM59t0WG377ouH3ZXr88C5_yhA4ak",
			},
		},
	}, nil
}
