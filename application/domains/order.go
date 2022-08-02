package domains

import (
	"errors"

	"github.com/jinzhu/gorm"
)

const (
	OrderStatus = "completed"
)

type (
	Order struct {
		gorm.Model
		Status     string
		Total      float64
		Country    string
		OrderItems []OrderItem
	}

	OrderItem struct {
		Qty       uint
		Price     float64
		Country   string
		Total     float64
		OrderID   uint
		ProductID uint
		Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (item *OrderItem) AddUnits(unit int) error {
	if unit <= 0 {
		return errors.New("Invalid units")
	}

	item.Qty += uint(unit)
	item.Total = item.Price * float64(item.Qty)
	return nil
}

func (order *Order) AddOrderItem(productId uint, price float64, qty int, country string) {
	newOrderItem := OrderItem{
		Qty:       uint(qty),
		Price:     price,
		Country:   country,
		Total:     price * float64(qty),
		ProductID: productId,
	}
	order.OrderItems = append(order.OrderItems, newOrderItem)
}

func (order Order) TotalValueOfOrderItems() int {
	total := 0
	for _, item := range order.OrderItems {
		total += int(item.Total)
	}
	return total
}
