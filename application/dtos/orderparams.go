package dtos

type (
	CreateOrderParam struct {
		Country    string           `json:"country" binding:"required"`
		OrderItems []OrderItemParam `json:"orderItems" binding:"required"`
	}

	OrderItemParam struct {
		Qty       uint    `json:"quantity" binding:"required"`
		Price     float64 `json:"price" binding:"required"`
		ProductId uint    `json:"productId" binding:"required"`
	}
)
