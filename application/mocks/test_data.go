package mock_services

import (
	"time"

	dto "github.com/philaden/mds-stock-keeping/application/dtos"
)

func GetTestProducts() []dto.ProductResponseDto {
	mockproducts := []dto.ProductResponseDto{
		{
			ID:             1,
			CreatedAt:      time.Now(),
			Name:           "Turner-Payne Soft Bike",
			Sku:            "da8ef851e075",
			Country:        "ke",
			AvailableStock: 100,
		},
		{
			ID:             2,
			CreatedAt:      time.Now(),
			Name:           "Turner-Payne Soft Bike",
			Sku:            "da8ef851e075",
			Country:        "ng",
			AvailableStock: 100,
		},
		{
			ID:             3,
			CreatedAt:      time.Now(),
			Name:           "Garcia, Jones and Murphy For repair Cotton Mouse",
			Sku:            "37d1fdc2cecb",
			Country:        "ci",
			AvailableStock: 100,
		},
	}

	return mockproducts
}
