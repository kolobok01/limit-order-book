package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/kolobublik/limit-order-book/api/orderbook/v1"
	"github.com/kolobublik/limit-order-book/internal/biz"
)

type OrderBookService struct {
	pb.UnimplementedOrderBookServer

	uc  *biz.OrderBookUsecase
	log *log.Helper
}

func NewOrderBookService(uc *biz.OrderBookUsecase, logger log.Logger) *OrderBookService {
	return &OrderBookService{uc: uc, log: log.NewHelper(logger)}
}

func (s *OrderBookService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	order, err := s.uc.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return toProtoOrder(order), nil
}

func (s *OrderBookService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return toProtoOrder(order), nil
}

func (s *OrderBookService) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	err := s.uc.CancelOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.CancelOrderResponse{OrderId: req.OrderId, Status: pb.OrderStatus_DONE}, nil
}

func (s *OrderBookService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.uc.ListOrders(ctx, req.Instrument)
	if err != nil {
		return nil, err
	}
	var protoOrders []*pb.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, toProtoOrder(order))
	}
	return &pb.ListOrdersResponse{Orders: protoOrders}, nil
}

func (s *OrderBookService) GetInstrumentQuote(ctx context.Context, req *pb.GetInstrumentQuoteRequest) (*pb.InstrumentQuote, error) {
	return s.uc.GetInstrumentQuote(ctx, req.Instrument)
}

func toProtoOrder(order *biz.Order) *pb.Order {
	return &pb.Order{
		OrderId:       order.ID,
		ClientOrderId: order.ClientOrderID,
		Instrument:    order.Instrument,
		Side:          order.Side,
		Type:          order.Type,
		Status:        order.Status,
		Price:         order.Price,
		Size:          order.Size,
		CreatedAt:     timestamppb.New(order.CreatedAt),
	}
}
