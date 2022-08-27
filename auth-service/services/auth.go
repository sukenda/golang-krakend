package services

import (
	"context"
	"encoding/base64"
	database "github.com/sukenda/golang-krakend/auth-service/database"
	models "github.com/sukenda/golang-krakend/auth-service/model"
	"github.com/sukenda/golang-krakend/auth-service/proto"
	"github.com/sukenda/golang-krakend/auth-service/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type AuthService struct {
	proto.UnimplementedAuthServiceServer
	Database database.Handler
	Jwt      utils.JwtWrapper
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
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &proto.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &proto.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *AuthService) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

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
		Alg: "RS256",
		Typ: "JWT",
		Kid: base64.StdEncoding.EncodeToString([]byte(s.Jwt.SecretKey)),
	}, nil
}
