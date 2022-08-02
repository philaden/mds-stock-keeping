package domain

import (
	"errors"
	"math"

	"github.com/jinzhu/gorm"
)

type (
	Product struct {
		gorm.Model
		Name           string
		Sku            string
		Country        string
		AvailableStock uint
	}
)

func (prd *Product) RemoveStock(quantityDesired int) (stockBalance uint, err error) {
	if prd.AvailableStock == 0 {
		return 0, errors.New("available stock balance can not be less than 0")
	}

	if quantityDesired > int(prd.AvailableStock) {
		return 0, errors.New("Invalid unit of quantity provided. Value can not be less than 0 or zero")
	}

	prd.AvailableStock -= uint(math.Abs(float64(quantityDesired)))
	return prd.AvailableStock, nil
}

func (prd *Product) AddStock(quantityDesired int) (stockBalance uint, err error) {

	originalBalance := prd.AvailableStock

	if quantityDesired <= 0 {
		return 0, errors.New("Invalid unit of quantity provided. Value can not be less than 0 or zero")
	}

	prd.AvailableStock += uint(quantityDesired)

	balance := prd.AvailableStock - originalBalance
	return uint(balance), nil
}

func GetNewStockValue(value int) uint {
	if ok := math.Signbit(float64(value)); ok {
		return 0
	}
	return uint(value)
}
