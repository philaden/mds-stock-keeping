package services

import (
	"errors"
	"fmt"

	domain "github.com/philaden/mds-stock-keeping/application/domains"
	dto "github.com/philaden/mds-stock-keeping/application/dtos"
	repo "github.com/philaden/mds-stock-keeping/application/repositories"
)

type (
	IOrderService interface {
		CreateOrder(model dto.CreateOrderParam) (uint, error)
	}

	OrderService struct {
		OrderRepository repo.IOrderRepository
	}
)

func NewOrderService(orderRepository repo.IOrderRepository) OrderService {
	return OrderService{OrderRepository: orderRepository}
}

func (orderService OrderService) CreateOrder(model dto.CreateOrderParam) (uint, error) {

	for _, item := range model.OrderItems {
		if item.Qty <= 0 {
			return 0, errors.New("invalid quantity of order item")
		}
	}

	order := domain.Order{Status: domain.OrderStatus, Country: model.Country}

	qty, available, err := orderService.OrderRepository.ValidateStock(model.Country, model.OrderItems)

	if err != nil {
		return qty, err
	}

	if !available {
		return qty, errors.New(fmt.Sprintf("Invalid order quantity, value is more than total stock quantity of %d", qty))
	}

	for _, item := range model.OrderItems {
		order.AddOrderItem(item.ProductId, item.Price, int(item.Qty), order.Country)
	}
	order.Total = float64(order.TotalValueOfOrderItems())

	if _, err := orderService.OrderRepository.CreateSingleOrder(order); err != nil {
		return 0, err
	}

	orderChannel := make(chan dto.OrderItemParam)
	go funnelInOrderItems(orderChannel, model.OrderItems)
	go orderService.OrderRepository.UpdateStockLevel(orderChannel, order.ID)

	return order.ID, nil
}

func funnelInOrderItems(orderChannel chan dto.OrderItemParam, model []dto.OrderItemParam) {
	for _, item := range model {
		orderChannel <- item
	}
}
