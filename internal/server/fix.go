package server

import (
	v1 "github.com/kolobublik/limit-order-book/api/fix/v1"
	"github.com/kolobublik/limit-order-book/internal/service"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func RegisterFIXService(srv *grpc.Server, fix *service.FIXService) {
	v1.RegisterFIXServiceServer(srv, fix)
}
