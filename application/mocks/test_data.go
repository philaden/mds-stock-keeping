package mocks

import (
	"time"

	domain "github.com/philaden/mds-stock-keeping/application/domains"
	dto "github.com/philaden/mds-stock-keeping/application/dtos"
)

func GetMockProductsDto() []dto.ProductResponseDto {
	mock := []dto.ProductResponseDto{
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

	return mock
}

func GetMockProducts() []domain.Product {
	mock := []domain.Product{
		{
			ID:             1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DeletedAt:      nil,
			Name:           "Turner-Payne Soft Bike",
			Sku:            "da8ef851e075",
			Country:        "ke",
			AvailableStock: 100,
		},
		{
			ID:             2,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DeletedAt:      nil,
			Name:           "Turner-Payne Soft Bike",
			Sku:            "da8ef851e075",
			Country:        "ng",
			AvailableStock: 200,
		},
		{
			ID:             3,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DeletedAt:      nil,
			Name:           "Garcia, Jones and Murphy For repair Cotton Mouse",
			Sku:            "37d1fdc2cecb",
			Country:        "ci",
			AvailableStock: 200,
		},
	}
	return mock
}

func GetMockProductBySku(sku string) *domain.Product {
	products := GetMockProducts()
	for _, prd := range products {
		if sku == prd.Sku {
			return &prd
		}
	}
	return nil
}

func CreateMockStockPayload() dto.UploadProductParam {
	return dto.UploadProductParam{
		Country:     "ci",
		Sku:         "37d1fdc2cecb",
		Name:        "Garcia, Jones and Murphy For repair Cotton Mouse",
		StockChange: 10,
	}
}

func CreateMockOrderPayload() domain.Order {
	return domain.Order{
		ID:      1,
		Status:  domain.OrderStatus,
		Total:   5000,
		Country: "ci",
		OrderItems: []domain.OrderItem{
			{
				OrderID: 1,
				ID:      1,
				Qty:     10,
				Price:   500,
				Country: "ci",
				Total:   5000,
			},
		},
	}
}
