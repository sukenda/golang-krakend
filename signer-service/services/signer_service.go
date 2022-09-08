package services

import (
	"context"
	"github.com/sukenda/golang-krakend/grpc-proto/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SignerService struct {
	proto.UnimplementedSignerServiceServer
}

func (s *SignerService) GetJsonWebKey(ctx context.Context, req *emptypb.Empty) (*proto.GetJsonWebKeyResponse, error) {
	return &proto.GetJsonWebKeyResponse{
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
