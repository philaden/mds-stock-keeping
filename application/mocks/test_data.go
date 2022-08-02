package mocks

import (
	domain "github.com/philaden/mds-stock-keeping/application/domains"
)

func GetTestProducts() []domain.Product {
	mockproducts := []domain.Product{
		{
			Name:           "Turner-Payne Soft Bike",
			Sku:            "da8ef851e075",
			Country:        "ke",
			AvailableStock: 100,
		},
		{
			Name:           "Turner-Payne Soft Bike",
			Sku:            "da8ef851e075",
			Country:        "ng",
			AvailableStock: 100,
		},
		{
			Name:           "Garcia, Jones and Murphy For repair Cotton Mouse",
			Sku:            "37d1fdc2cecb",
			Country:        "ci",
			AvailableStock: 100,
		},
	}

	return mockproducts
}
