package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	pb "github.com/kolobublik/limit-order-book/api/orderbook/v1"
)

// Order is a Order model.
// OrderBookUsecase is a Greeter usecase.
type OrderBookUsecase struct {
	repo OrderBookRepo
	log  *log.Helper
}

// NewOrderBookUsecase new a Greeter usecase.
func NewOrderBookUsecase(repo OrderBookRepo, logger log.Logger) *OrderBookUsecase {
	return &OrderBookUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateOrder creates a Order, and returns the new Order.
func (uc *OrderBookUsecase) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*Order, error) {
	uc.log.WithContext(ctx).Infof("CreateOrder: %v", req)
	order := &Order{
		ClientOrderID: req.ClientOrderId,
		Instrument:    req.Instrument,
		Side:          req.Side,
		Type:          req.Type,
		Price:         req.Price,
		Size:          req.Size,
		ID:            uuid.New().String(),
	}
	return uc.repo.CreateOrder(ctx, order)
}

// GetOrder gets a Order.
func (uc *OrderBookUsecase) GetOrder(ctx context.Context, id string) (*Order, error) {
	uc.log.WithContext(ctx).Infof("GetOrder: %v", id)
	return uc.repo.GetOrder(ctx, id)
}

// CancelOrder cancels a Order.
func (uc *OrderBookUsecase) CancelOrder(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("CancelOrder: %v", id)
	return uc.repo.CancelOrder(ctx, id)
}

// ListOrders lists Orders.
func (uc *OrderBookUsecase) ListOrders(ctx context.Context, instrument string) ([]*Order, error) {
	uc.log.WithContext(ctx).Infof("ListOrders: %v", instrument)
	return uc.repo.ListOrders(ctx)
}

// GetInstrumentQuote gets an InstrumentQuote.
func (uc *OrderBookUsecase) GetInstrumentQuote(ctx context.Context, instrument string) (*pb.InstrumentQuote, error) {
	uc.log.WithContext(ctx).Infof("GetInstrumentQuote: %v", instrument)
	// Convert ticker to instrument quote
	ticker, err := uc.repo.GetTicker(ctx, instrument)
	if err != nil {
		return nil, err
	}
	return &pb.InstrumentQuote{
		Instrument:   ticker.ProductID,
		BestBidPrice: ticker.Price,
		BestBidSize:  ticker.Volume,
		BestAskPrice: ticker.Price,
		BestAskSize:  ticker.Volume,
	}, nil
}
