package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1fix "github.com/kolobublik/limit-order-book/api/fix/v1"
	v1ob "github.com/kolobublik/limit-order-book/api/orderbook/v1"
	"github.com/kolobublik/limit-order-book/internal/conf"
	"github.com/kolobublik/limit-order-book/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, orderbook *service.OrderBookService, fix *service.FIXService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1ob.RegisterOrderBookServer(srv, orderbook)
	v1fix.RegisterFIXServiceServer(srv, fix)
	return srv
}
