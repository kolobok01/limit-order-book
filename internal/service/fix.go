package service

import (
	"context"

	v1 "github.com/kolobublik/limit-order-book/api/fix/v1"
	"github.com/kolobublik/limit-order-book/internal/data/fix/model"
)

func (s *FIXService) PlaceOrder(ctx context.Context, req *v1.PlaceOrderRequest) (*v1.PlaceOrderResponse, error) {
	// Convert proto enum to string
	var orderType, side, qtyType string

	switch req.Order.OrderType {
	case v1.OrderType_ORDER_TYPE_MARKET:
		orderType = "MARKET"
	case v1.OrderType_ORDER_TYPE_LIMIT:
		orderType = "LIMIT"
	case v1.OrderType_ORDER_TYPE_VWAP:
		orderType = "VWAP"
	}

	switch req.Order.Side {
	case v1.Side_SIDE_BUY:
		side = "BUY"
	case v1.Side_SIDE_SELL:
		side = "SELL"
	}

	switch req.Order.QuantityType {
	case v1.QuantityType_QUANTITY_TYPE_BASE:
		qtyType = "BASE"
	case v1.QuantityType_QUANTITY_TYPE_QUOTE:
		qtyType = "QUOTE"
	}

	// Create optional VWAP parameters array
	var vwapParams []string
	if req.Order.OrderType == v1.OrderType_ORDER_TYPE_VWAP {
		vwapParams = []string{
			req.Order.StartTime,
			req.Order.ParticipationRate,
			req.Order.ExpireTime,
		}
	}

	// Place the order using the FIX client
	err := s.client.PlaceOrder(
		req.Order.Symbol,
		orderType,
		side,
		qtyType,
		req.Order.Quantity,
		req.Order.Price,
		"", // portfolio ID will be taken from client config
		vwapParams...,
	)

	if err != nil {
		return nil, err
	}

	return &v1.PlaceOrderResponse{
		Status: "PLACED",
	}, nil
}

func (s *FIXService) GetOrderStatus(ctx context.Context, req *v1.GetOrderStatusRequest) (*v1.GetOrderStatusResponse, error) {
	info := model.OrderInfo{
		ClOrdId:  req.OrderInfo.ClOrdId,
		OrderId:  req.OrderInfo.OrderId,
		Symbol:   req.OrderInfo.Symbol,
		Side:     req.OrderInfo.Side,
		Status:   req.OrderInfo.Status,
		Type:     req.OrderInfo.Type,
		Price:    req.OrderInfo.Price,
		Quantity: req.OrderInfo.Quantity,
	}

	err := s.client.GetOrderStatus(info, "") // portfolio ID will be taken from client config
	if err != nil {
		return nil, err
	}

	return &v1.GetOrderStatusResponse{
		OrderInfo: req.OrderInfo,
	}, nil
}

func (s *FIXService) CancelOrder(ctx context.Context, req *v1.CancelOrderRequest) (*v1.CancelOrderResponse, error) {
	info := model.OrderInfo{
		ClOrdId:  req.OrderInfo.ClOrdId,
		OrderId:  req.OrderInfo.OrderId,
		Symbol:   req.OrderInfo.Symbol,
		Side:     req.OrderInfo.Side,
		Status:   req.OrderInfo.Status,
		Type:     req.OrderInfo.Type,
		Price:    req.OrderInfo.Price,
		Quantity: req.OrderInfo.Quantity,
	}

	err := s.client.CancelOrder(info, "") // portfolio ID will be taken from client config
	if err != nil {
		return nil, err
	}

	return &v1.CancelOrderResponse{
		Status: "CANCELLED",
	}, nil
}
