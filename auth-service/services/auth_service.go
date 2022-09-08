package services

import (
	"context"
	"errors"
	database "github.com/sukenda/golang-krakend/auth-service/database"
	models "github.com/sukenda/golang-krakend/auth-service/model"
	"github.com/sukenda/golang-krakend/auth-service/utils"
	"github.com/sukenda/golang-krakend/proto"
	"google.golang.org/protobuf/types/known/emptypb"
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
		Aud: " http://auth-service:8081", // Auth service URL
		Exp: float32(exp),
		Sub: "1234567890qwertyuio",
		Jti: "mnb23vcsrt756yuiomnbvcx98ertyuiop",
		Iss: "http://localhost:8080", // Krakend URL
		Claims: &proto.Claims{
			UserId: user.ID.String(),
			Email:  user.Email,
			Roles:  []string{"role_a", "role_b"},
		},
	}

	refresh := &proto.Token{
		Sub: "1234567890qwertyuio",
		Jti: "mnb23vcsrt756yuiomnbvcx98ertyuiop",
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

func (s *AuthService) JWKValidate(ctx context.Context, req *emptypb.Empty) (*proto.JWKValidateResponse, error) {
	return &proto.JWKValidateResponse{
		Keys: []*proto.JWK{
			{
				Kty: "oct",
				Use: "sig",
				Kid: "bluebird.id",
				K:   "pWn7Tu6Jz8EQ4eHFiGVgmroA4_ENLvqLAUnMxxPx4epMpRNQNtPp86DHBq-kU5Es4V5rk4O6cCD1pCS1-IMy_I_w9yeA5o6-AnK4iMSiXLa9_9RAygO3Xb2NMhlI6CDduTA85nhRbm8TCLOKZTTX2QRAn3yoGY1arw1HrST-FDusWjOmIfGggMa2GZ9MD1y1v0XFix7ACRyEqS7EgSOBgLo2HOJYEE-ZHrULUNhzCG8CljD7AyYPo3iaxZJcmDLZzoSsAsJMULcx5rQmFNjUGMlyGcsLDHklWj4UFKATKP3tZPBvxAJpKzUyFdJYGKzg8IUY6ZhLGqpEr09RcWbPpg",
				Alg: "HS256",
			},
		},
	}, nil
}
