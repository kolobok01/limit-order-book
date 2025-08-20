package data

import (
	"context"
	"fmt"
	"strconv"

	"github.com/coinbase-samples/advanced-trade-sdk-go/model"
	"github.com/coinbase-samples/advanced-trade-sdk-go/orders"
	"github.com/coinbase-samples/advanced-trade-sdk-go/products"
	"github.com/go-kratos/kratos/v2/log"

	pb "github.com/kolobublik/limit-order-book/api/orderbook/v1"
	"github.com/kolobublik/limit-order-book/internal/biz"
)

type orderBookRepo struct {
	data            *Data
	log             *log.Helper
	ordersService   orders.OrdersService
	productsService products.ProductsService
}

func NewOrderBookRepo(data *Data, logger log.Logger) biz.OrderBookRepo {
	return &orderBookRepo{
		data:            data,
		log:             log.NewHelper(logger),
		ordersService:   orders.NewOrdersService(data.CoinbaseClient),
		productsService: products.NewProductsService(data.CoinbaseClient),
	}
}

func (r *orderBookRepo) CreateOrder(ctx context.Context, order *biz.Order) (*biz.Order, error) {
	orderRequest := &orders.CreateOrderRequest{
		ClientOrderId: order.ClientOrderID,
		ProductId:     order.Instrument,
		Side:          order.Side.String(),
		OrderConfiguration: model.OrderConfiguration{
			MarketMarketIoc: &model.MarketIoc{
				BaseSize: fmt.Sprintf("%f", order.Size),
			},
		},
	}

	resp, err := r.ordersService.CreateOrder(ctx, orderRequest)
	if err != nil {
		return nil, err
	}

	if resp.SuccessResponse == nil {
		return nil, fmt.Errorf("order creation failed: no success response")
	}

	return &biz.Order{
		ID:            resp.SuccessResponse.OrderId,
		ClientOrderID: resp.SuccessResponse.ClientOrderId,
		Instrument:    resp.SuccessResponse.ProductId,
		Type:          pb.OrderType_MARKET,
		Side:          mapOrderSide(resp.SuccessResponse.Side),
		Size:          order.Size,
		Price:         order.Price,
		Status:        pb.OrderStatus_WORKING,
	}, nil
}

func (r *orderBookRepo) GetOrder(ctx context.Context, id string) (*biz.Order, error) {
	getOrderRequest := &orders.GetOrderRequest{
		OrderId: id,
	}

	resp, err := r.ordersService.GetOrder(ctx, getOrderRequest)
	if err != nil {
		return nil, err
	}

	if resp.Order == nil {
		return nil, fmt.Errorf("order not found: %s", id)
	}

	return &biz.Order{
		ID:            resp.Order.OrderId,
		ClientOrderID: resp.Order.ClientOrderId,
		Instrument:    resp.Order.ProductId,
		Type:          mapOrderType(resp.Order.OrderConfiguration),
		Side:          mapOrderSide(resp.Order.Side),
		Size:          parseFloat(resp.Order.FilledSize),
		Price:         parseFloat(resp.Order.AverageFilledPrice),
		Status:        mapOrderStatus(resp.Order.Status),
	}, nil
}

func (r *orderBookRepo) CancelOrder(ctx context.Context, id string) error {
	cancelOrderRequest := &orders.CancelOrdersRequest{
		OrderIds: []string{id},
	}

	resp, err := r.ordersService.CancelOrders(ctx, cancelOrderRequest)
	if err != nil {
		return err
	}

	if len(resp.Results) > 0 && resp.Results[0].Success {
		return nil
	}

	return fmt.Errorf("failed to cancel order: %s", id)
}

func (r *orderBookRepo) ListOrders(ctx context.Context) ([]*biz.Order, error) {
	resp, err := r.ordersService.ListOrders(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []*biz.Order
	for _, order := range resp.Orders {
		result = append(result, &biz.Order{
			ID:            order.OrderId,
			ClientOrderID: order.ClientOrderId,
			Instrument:    order.ProductId,
			Type:          mapOrderType(order.OrderConfiguration),
			Side:          mapOrderSide(order.Side),
			Size:          parseFloat(order.FilledSize),
			Price:         parseFloat(order.AverageFilledPrice),
			Status:        mapOrderStatus(order.Status),
		})
	}

	return result, nil
}

func (r *orderBookRepo) GetTicker(ctx context.Context, productID string) (*biz.Ticker, error) {
	getProductRequest := &products.GetProductRequest{
		ProductId: productID,
	}

	resp, err := r.productsService.GetProduct(ctx, getProductRequest)
	if err != nil {
		return nil, err
	}

	return &biz.Ticker{
		ProductID: resp.ProductId,
		Price:     parseFloat(resp.Price),
		Volume:    parseFloat(resp.Volume24h),
		Low:       0, // TODO: handle low/high/open once historical data is available
		High:      0,
		Open:      0,
	}, nil
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func mapOrderType(config model.OrderConfiguration) pb.OrderType {
	switch {
	case config.MarketMarketIoc != nil:
		return pb.OrderType_MARKET
	case config.LimitLimitGtc != nil:
		return pb.OrderType_LIMIT
	default:
		return pb.OrderType_ORDER_TYPE_UNSPECIFIED
	}
}

func mapOrderStatus(status string) pb.OrderStatus {
	switch status {
	case "OPEN":
		return pb.OrderStatus_WORKING
	case "FILLED", "CANCELLED":
		return pb.OrderStatus_DONE
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func mapOrderSide(side string) pb.Side {
	switch side {
	case "BUY":
		return pb.Side_BUY
	case "SELL":
		return pb.Side_SELL
	default:
		return pb.Side_SIDE_UNSPECIFIED
	}
}
