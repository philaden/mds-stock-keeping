package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philaden/mds-stock-keeping/application/params"
	Inject "github.com/philaden/mds-stock-keeping/application/services"
	"github.com/philaden/mds-stock-keeping/infrastructure"
)

type ApiService struct {
	Product_Service Inject.IProductService
	order_Service   Inject.IOrderService
}

func SetupContollerRoutes(router *gin.Engine) {

	services := ApiService{
		Product_Service: Inject.ProductService{
			DbContext: infrastructure.Connection,
		},

		order_Service: Inject.OrderService{
			DbContext: infrastructure.Connection,
		},
	}

	controllerRouter := router.Group("/api")

	controllerRouter.POST("/products", services.HandleStockBulkUpload)
	controllerRouter.GET("/products/:sku", services.HandleGetProductBySku)
	controllerRouter.GET("/products", services.HandleGetProducts)
	controllerRouter.POST("/products/single", services.HandleSingleStock)
	controllerRouter.POST("/orders", services.HandleCreateOrder)
}

// @Summary update stocks
// @Produce json
// @Description This endpoint accept stock information as a csv file
// @Router /api/products [post]
func (apiServices ApiService) HandleStockBulkUpload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Operation failed", "status": false, "error": err.Error()})
		fmt.Printf(err.Error())
		return
	}

	go apiServices.Product_Service.UploadStock(file)

	c.JSON(http.StatusOK, gin.H{
		"message": "stock uploaded successfully",
		"data":    nil,
		"status":  true,
	})
}

// @Summary Attempts to get a existing product by sku
// @Produce json
// @Description This endpoint fetches a stock product by its sku
// @Router /api/products/:sku [get]
func (apiServices ApiService) HandleGetProductBySku(c *gin.Context) {

	skuParam, ok := c.GetQuery("sku")

	if ok && skuParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No sku name found, please try again."})
		c.Abort()
		return
	}

	response, err := apiServices.Product_Service.GetProductBySku(skuParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Operation Successful",
		"data":    response,
		"status":  true,
	})

}

// @Summary Attempts to get all products
// @Produce json
// @Description This endpoint fetches a list of all create product
// @Router /api/products [get]
func (apiServices ApiService) HandleGetProducts(c *gin.Context) {

	response, err := apiServices.Product_Service.GetProducts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Operation Successful",
		"data":    response,
		"status":  true,
	})

}

// @Summary Create a new single product
// @Produce json
// @Description This endpoint creates a single stock product
// @Router /api/products/single [post]
func (apiServices ApiService) HandleSingleStock(c *gin.Context) {

	var json params.UploadProductParam

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect details supplied, please try again."})
		return
	}

	status, err := apiServices.Product_Service.CreateSingleStock(json)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The product has been successfully created.",
		"data":    nil,
		"status":  status,
	})
}

// @Summary Create a new order
// @Description This endpoint creates an order for a product. This endpoint is meant to reduce the unit of stock for product
// @Produce json
// @Router /order [post]
func (apiServices ApiService) HandleCreateOrder(c *gin.Context) {

	var json params.CreateOrderParam

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect details supplied, please try again."})
		return
	}

	insertedId, err := apiServices.order_Service.CreateOrder(json)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your order has been successfully created.",
		"orderId": insertedId,
		"status":  true,
	})
}
