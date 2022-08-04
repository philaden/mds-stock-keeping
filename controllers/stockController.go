package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/philaden/mds-stock-keeping/application/dtos"
	"github.com/philaden/mds-stock-keeping/application/repositories"
	"github.com/philaden/mds-stock-keeping/application/services"
	"github.com/philaden/mds-stock-keeping/infrastructure"
)

type ApiService struct {
	ProductService services.IProductService
	OrderService   services.IOrderService
}

func NewApiService() ApiService {
	productRepo := repositories.NewProductRepostiory(infrastructure.Connection)
	orderRepo := repositories.NewOrderRepostiory(infrastructure.Connection)

	return ApiService{
		ProductService: services.NewProductService(productRepo),
		OrderService:   services.NewOrderService(orderRepo),
	}
}

func SetupContollerRoutes(router *gin.Engine) {
	services := NewApiService()
	controller := router.Group("/api")

	controller.POST("/products", services.HandleStockBulkUpload)
	controller.GET("/products/:sku", services.HandleGetProductBySku)
	controller.GET("/products", services.HandleGetProducts)
	controller.POST("/products/single", services.HandleSingleStock)
	controller.POST("/orders", services.HandleCreateOrder)
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

	fileType := strings.Split(file.Header.Get("Content-Type"), "/")[1]
	if fileType != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file type. Please upload a csv file",
			"status":  false,
			"error":   errors.New("Invalid file type. Please upload a csv file"),
		})
	}

	go apiServices.ProductService.UploadStock(file)

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

	response, err := apiServices.ProductService.GetProductBySku(skuParam)

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

	response, err := apiServices.ProductService.GetProducts()

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

	var json dto.UploadProductParam

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect details supplied, please try again."})
		return
	}

	status, err := apiServices.ProductService.CreateSingleStock(json)

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

	var json dto.CreateOrderParam

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect details supplied, please try again."})
		return
	}

	insertedId, err := apiServices.OrderService.CreateOrder(json)

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
