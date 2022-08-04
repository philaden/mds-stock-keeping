package repositories

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	domain "github.com/philaden/mds-stock-keeping/application/domains"
	dto "github.com/philaden/mds-stock-keeping/application/dtos"
)

type (
	IOrderRepository interface {
		CreateSingleOrder(model domain.Order) (uint, error)
		UpdateStockLevel(orderChannel chan dto.OrderItemParam, orderId uint) error
		GetOrderItemsByProductIds(productIds []uint) (orderedItems []domain.OrderItem)
		ValidateStock(country string, orderItems []dto.OrderItemParam) (stockBalance uint, isAvailable bool, err error)
	}

	OrderRepository struct {
		DbContext *gorm.DB
	}
)

func NewOrderRepostiory(dbContext *gorm.DB) IOrderRepository {
	return OrderRepository{DbContext: dbContext}
}

func (repo OrderRepository) CreateSingleOrder(model domain.Order) (uint, error) {
	if err := repo.DbContext.Create(&model).Error; err != nil {
		return 0, err
	}
	return model.ID, nil
}

func (repo OrderRepository) ValidateStock(country string, orderItems []dto.OrderItemParam) (stockBalance uint, isAvailable bool, err error) {
	for _, orderItem := range orderItems {
		var product domain.Product

		if err := repo.DbContext.First(&product, orderItem.ProductId).Error; err != nil {
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

func (repo OrderRepository) UpdateStockLevel(orderChannel chan dto.OrderItemParam, orderId uint) error {

	var order domain.Order

	if err := repo.DbContext.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		fmt.Printf(err.Error())
		return err
	}

	var orderedItems []domain.OrderItem
	productIds := make([]uint, len(order.OrderItems)-1)

	for _, item := range order.OrderItems {
		productIds = append(productIds, item.ProductID)
	}

	orderedItems = repo.GetOrderItemsByProductIds(productIds)

	for {
		item := <-orderChannel
		for _, oItem := range orderedItems {
			if oItem.ProductID == item.ProductId && oItem.OrderID == orderId {
				var prod domain.Product
				if err := repo.DbContext.First(&prod, oItem.ProductID).Error; err == nil {
					prod.AvailableStock = prod.AvailableStock - item.Qty
					repo.DbContext.Save(&prod)
				}
			}
		}
	}
}

func (repo OrderRepository) GetOrderItemsByProductIds(productIds []uint) (orderedItems []domain.OrderItem) {

	for _, pId := range productIds {
		var subQueryItems []domain.OrderItem
		repo.DbContext.Model("order_items").Where(&domain.OrderItem{ProductID: pId}).Find(&subQueryItems)
		for _, item := range subQueryItems {
			orderedItems = append(orderedItems, item)
		}
	}
	return orderedItems
}
