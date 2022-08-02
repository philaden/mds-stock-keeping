package services

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	domain "github.com/philaden/mds-stock-keeping/application/domain"
	"github.com/philaden/mds-stock-keeping/application/params"
)

type (
	IOrderService interface {
		CreateOrder(model params.CreateOrderParam) (uint, error)
	}

	OrderService struct {
		DbContext *gorm.DB
	}
)

func (orderService OrderService) CreateOrder(model params.CreateOrderParam) (uint, error) {

	for _, item := range model.OrderItems {
		if item.Qty <= 0 {
			return 0, errors.New("invalid quantity of order item")
		}
	}

	order := domain.Order{Status: domain.OrderStatus, Country: model.Country}

	qty, available, err := orderService.validateStock(model.Country, model.OrderItems)

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

	if err := orderService.DbContext.Create(&order).Error; err != nil {
		return 0, err
	}

	orderChannel := make(chan params.OrderItemParam)
	go funnelInOrderItems(orderChannel, model.OrderItems)
	go orderService.updateStockLevel(orderChannel, order.ID)

	return order.ID, nil
}

func (orderService OrderService) validateStock(country string, orderItems []params.OrderItemParam) (stockBalance uint, isAvailable bool, err error) {
	for _, orderItem := range orderItems {
		var product domain.Product

		if err := orderService.DbContext.First(&product, orderItem.ProductId).Error; err != nil {
			return 0, false, err
		}

		stockBalance += product.AvailableStock

		if product.Country != country {
			return 0, false, errors.New("This product does not belong to the country you are ordering from")
		}

		if orderItem.Qty > product.AvailableStock {
			return 0, false, errors.New("ordered quantity is greater than the value of stock")
		}
	}
	isAvailable = true
	return stockBalance, isAvailable, nil
}

func funnelInOrderItems(orderChannel chan params.OrderItemParam, model []params.OrderItemParam) {
	for _, item := range model {
		orderChannel <- item
	}
}

func (orderService OrderService) updateStockLevel(orderChannel chan params.OrderItemParam, orderId uint) error {

	var order domain.Order

	if err := orderService.DbContext.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		fmt.Printf(err.Error())
		return err
	}

	var orderedItems []domain.OrderItem
	productIds := make([]uint, len(order.OrderItems)-1)

	for _, item := range order.OrderItems {
		productIds = append(productIds, item.ProductID)
	}

	orderedItems = orderService.getOrderItemsByProductIds(productIds)

	for {
		item := <-orderChannel
		for _, oItem := range orderedItems {
			if oItem.ProductID == item.ProductId && oItem.OrderID == orderId {
				var prod domain.Product
				if err := orderService.DbContext.First(&prod, oItem.ProductID).Error; err == nil {
					prod.AvailableStock = prod.AvailableStock - item.Qty
					orderService.DbContext.Save(&prod)
				}
			}
		}
	}
}

func (orderService OrderService) getOrderItemsByProductIds(productIds []uint) (orderedItems []domain.OrderItem) {

	for _, pId := range productIds {
		var subQueryItems []domain.OrderItem
		orderService.DbContext.Model("order_items").Where(&domain.OrderItem{ProductID: pId}).Find(&subQueryItems)
		for _, item := range subQueryItems {
			orderedItems = append(orderedItems, item)
		}
	}
	return orderedItems
}
