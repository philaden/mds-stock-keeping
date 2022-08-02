package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"strconv"

	"github.com/jinzhu/gorm"
	domain "github.com/philaden/mds-stock-keeping/application/domains"
	params "github.com/philaden/mds-stock-keeping/application/dtos"
)

type (
	IProductService interface {
		UploadStock(file *multipart.FileHeader) (bool, error)
		GetProducts() (products []domain.Product, err error)
		GetProductBySku(sku string) (prd *domain.Product, err error)
		CreateSingleStock(stock params.UploadProductParam) (bool, error)
	}

	ProductService struct {
		DbContext *gorm.DB
	}
)

func (productService ProductService) UploadStock(file *multipart.FileHeader) (bool, error) {

	openedFile, err := file.Open()
	defer openedFile.Close()

	if err != nil {
		fmt.Printf(err.Error())
		return false, err
	}

	csvReader := csv.NewReader(openedFile)
	lines, err := csvReader.ReadAll()

	if err != nil {
		return false, err
	}

	var stocks []params.UploadProductParam
	for i := 0; i < len(lines); i++ {
		if i == 0 {
			continue
		}

		line := lines[i]
		stockChange, err := strconv.Atoi(line[3])

		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		uploadData := params.UploadProductParam{Country: line[0], Sku: line[1], Name: line[2], StockChange: stockChange}
		stocks = append(stocks, uploadData)
	}
	return productService.SaveStocks(stocks)
}

func (productService ProductService) SaveStocks(stocks []params.UploadProductParam) (bool, error) {
	for _, stock := range stocks {
		var prd *domain.Product = &domain.Product{}

		if err := productService.DbContext.Where(&domain.Product{Sku: stock.Sku, Country: stock.Country}).First(&prd).Error; err == nil {
			if ok := math.Signbit(float64(stock.StockChange)); ok {
				if _, err := prd.RemoveStock(stock.StockChange); err != nil {
					return false, err
				}
			} else {
				if _, err := prd.AddStock(stock.StockChange); err != nil {
					return false, err
				}
			}

			if err := productService.DbContext.Save(&prd).Error; err != nil {
				return false, err
			}
		} else {
			newProduct := domain.Product{
				Name:           stock.Name,
				Sku:            stock.Sku,
				Country:        stock.Country,
				AvailableStock: domain.GetNewStockValue(stock.StockChange),
			}

			if err := productService.DbContext.Create(&newProduct).Error; err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (productService ProductService) GetProducts() (products []domain.Product, err error) {
	if err := productService.DbContext.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (productService ProductService) GetProductBySku(sku string) (prd *domain.Product, err error) {
	prd = &domain.Product{}
	if err := productService.DbContext.Where(&domain.Product{Sku: sku}).First(&prd).Error; err != nil {
		return nil, err
	}
	return prd, nil
}

func (productService ProductService) CreateSingleStock(stock params.UploadProductParam) (bool, error) {
	var prd *domain.Product = &domain.Product{}
	payload, _ := json.Marshal(stock)
	fmt.Printf(string(payload))

	if err := productService.DbContext.Where(&domain.Product{Sku: stock.Sku}).First(&prd).Error; err == nil {
		if ok := math.Signbit(float64(stock.StockChange)); ok {
			if _, err := prd.RemoveStock(stock.StockChange); err != nil {
				return false, err
			}
		} else {
			if _, err := prd.AddStock(stock.StockChange); err != nil {
				return false, err
			}
		}
		if err := productService.DbContext.Save(&prd).Error; err != nil {
			return false, err
		}
	} else {
		newProduct := domain.Product{
			Name:           stock.Name,
			Sku:            stock.Sku,
			Country:        stock.Country,
			AvailableStock: 0,
		}

		if err := productService.DbContext.Create(&newProduct).Error; err != nil {
			return false, err
		}
	}
	return true, nil
}
