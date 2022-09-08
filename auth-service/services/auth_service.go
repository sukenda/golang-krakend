package services

import (
	"context"
	"errors"
	database "github.com/sukenda/golang-krakend/auth-service/database"
	models "github.com/sukenda/golang-krakend/auth-service/model"
	"github.com/sukenda/golang-krakend/auth-service/utils"
	"github.com/sukenda/golang-krakend/grpc-proto/proto"
	"net/http"
	"time"
)

type AuthService struct {
	proto.UnimplementedAuthServiceServer
	Database   database.Handler
	JwtWrapper utils.JwtWrapper
	JWT        utils.JWT
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
		return nil, errors.New("user not found")
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		return nil, errors.New("password not match")
	}

	exp := time.Now().Local().Add(time.Hour * time.Duration(10)).Unix()
	access := &proto.Token{
		Sub: user.ID.String(),
		Jti: user.ID.String(),
		Aud: " http://auth-service:8081", // Auth service URL
		Exp: float32(exp),
		Iss: "http://localhost:8080", // Krakend URL
		Claims: &proto.Claims{
			UserId: user.ID.String(),
			Email:  user.Email,
			Roles:  []string{"role_a", "role_b"},
		},
	}

	refresh := &proto.Token{
		Sub: user.ID.String(),
		Jti: user.ID.String(),
		Aud: " http://auth-service:8081", // Auth service URL
		Exp: float32(exp),
		Iss: "http://localhost:8080", // Krakend URL
		Claims: &proto.Claims{
			UserId: user.ID.String(),
			Email:  user.Email,
		},
	}

	return &proto.LoginResponse{
		Exp:          float32(exp),
		AccessToken:  access,
		RefreshToken: refresh,
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
		UserId: user.ID.String(),
	}, nil
}
