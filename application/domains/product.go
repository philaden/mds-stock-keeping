package domains

import (
	"errors"
	"math"
	"time"

	dto "github.com/philaden/mds-stock-keeping/application/dtos"
)

type (
	Product struct {
		ID             uint `gorm:"primaryKey"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		DeletedAt      *time.Time `gorm:"index"`
		Name           string
		Sku            string
		Country        string
		AvailableStock uint
	}

	Products []Product
)

func ToDto(prd *Product) *dto.ProductResponseDto {
	return &dto.ProductResponseDto{
		ID:             prd.ID,
		CreatedAt:      prd.CreatedAt,
		Name:           prd.Name,
		Sku:            prd.Sku,
		Country:        prd.Country,
		AvailableStock: prd.AvailableStock,
	}
}

func ToSliceDto(prds Products) dto.ProductsResponseDto {
	var dtos []dto.ProductResponseDto
	for _, value := range prds {
		dtos = append(dtos, dto.ProductResponseDto{
			ID:             value.ID,
			CreatedAt:      value.CreatedAt,
			Name:           value.Name,
			Sku:            value.Sku,
			Country:        value.Country,
			AvailableStock: value.AvailableStock,
		})
	}
	return dtos
}

func (prd *Product) RemoveStock(stockChange int) (stockBalance uint, err error) {
	if prd.AvailableStock == 0 {
		return 0, errors.New("available stock balance can not be less than 0")
	}

	if stockChange > int(prd.AvailableStock) {
		return 0, errors.New("Invalid unit of quantity provided. Value can not be less than 0 or zero")
	}

	prd.AvailableStock -= uint(stockChange)
	return prd.AvailableStock, nil
}

func (prd *Product) AddStock(stockChange int) (uint, error) {

	if stockChange <= 0 {
		return 0, errors.New("Invalid unit of quantity provided. Value can not be less than 0 or zero")
	}

	prd.AvailableStock += uint(stockChange)

	return uint(prd.AvailableStock), nil
}

func GetNewStockValue(value int) uint {
	if ok := math.Signbit(float64(value)); ok {
		return 0
	}
	return uint(value)
}
