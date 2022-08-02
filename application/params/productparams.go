package params

type (
	UploadProductParam struct {
		Country     string `json:"country" binding:"required"`
		Sku         string `json:"sku" binding:"required"`
		Name        string `json:"name" binding:"required"`
		StockChange int    `json:"stockchange" binding:"required"`
	}
)
