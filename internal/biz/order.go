package biz

import (
	"context"
	"time"

	pb "github.com/kolobublik/limit-order-book/api/orderbook/v1"
)

type Order struct {
	ID            string
	ClientOrderID string
	Instrument    string
	Side          pb.Side
	Type          pb.OrderType
	Status        pb.OrderStatus
	Price         float64
	Size          float64
	CreatedAt     time.Time
}

type Ticker struct {
	ProductID string
	Price     float64
	Volume    float64
	Low       float64
	High      float64
	Open      float64
}

type OrderBookRepo interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	CancelOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context) ([]*Order, error)
	GetTicker(ctx context.Context, productID string) (*Ticker, error)
}
