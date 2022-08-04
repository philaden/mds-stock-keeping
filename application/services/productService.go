package services

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"

	domain "github.com/philaden/mds-stock-keeping/application/domains"
	dto "github.com/philaden/mds-stock-keeping/application/dtos"
	repo "github.com/philaden/mds-stock-keeping/application/repositories"
)

type (
	IProductService interface {
		UploadStock(file *multipart.FileHeader) (bool, error)
		GetProducts() (dto.ProductsResponseDto, error)
		GetProductBySku(sku string) (*dto.ProductResponseDto, error)
		CreateSingleStock(stock dto.UploadProductParam) (bool, error)
	}

	ProductService struct {
		ProductRepository repo.IProductRepository
	}
)

func NewProductService(productRepository repo.IProductRepository) ProductService {
	return ProductService{ProductRepository: productRepository}
}

func (service ProductService) UploadStock(file *multipart.FileHeader) (bool, error) {

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

	var stocks []dto.UploadProductParam
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

		uploadData := dto.UploadProductParam{Country: line[0], Sku: line[1], Name: line[2], StockChange: stockChange}
		stocks = append(stocks, uploadData)
	}
	return service.SaveStocks(stocks)
}

func (service ProductService) SaveStocks(stocks []dto.UploadProductParam) (bool, error) {
	for _, stock := range stocks {
		service.ProductRepository.SaveStock(stock.Country, stock.Sku, stock.Name, stock.StockChange)
	}
	return true, nil
}

func (service ProductService) GetProducts() (dto.ProductsResponseDto, error) {
	data, err := service.ProductRepository.GetProducts()
	if err != nil {
		return nil, err
	}
	return domain.ToSliceDto(data), nil
}

func (service ProductService) GetProductBySku(sku string) (*dto.ProductResponseDto, error) {
	data, err := service.ProductRepository.GetProductBySku(sku)
	if err != nil {
		return nil, err
	}
	return domain.ToDto(data), nil
}

func (service ProductService) CreateSingleStock(stock dto.UploadProductParam) (bool, error) {
	return service.ProductRepository.SaveStock(stock.Country, stock.Name, stock.Sku, stock.StockChange)
}
