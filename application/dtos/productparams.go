package dtos

import (
	"time"
)

type (
	UploadProductParam struct {
		Country     string `json:"country" binding:"required"`
		Sku         string `json:"sku" binding:"required"`
		Name        string `json:"name" binding:"required"`
		StockChange int    `json:"stockchange" binding:"required"`
	}

	ProductResponseDto struct {
		ID             uint `gorm:"primaryKey"`
		CreatedAt      time.Time
		Name           string
		Sku            string
		Country        string
		AvailableStock uint
	}

	ProductsResponseDto []ProductResponseDto
)
