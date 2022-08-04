package repositories

import (
	"math"

	"github.com/jinzhu/gorm"
	domain "github.com/philaden/mds-stock-keeping/application/domains"
)

type (
	IProductRepository interface {
		GetProducts() ([]domain.Product, error)
		GetProductBySku(sku string) (*domain.Product, error)
		SaveStock(country, sku, name string, stockChange int) (bool, error)
	}

	ProductRepository struct {
		DbContext *gorm.DB
	}
)

func NewProductRepostiory(dbContext *gorm.DB) IProductRepository {
	return ProductRepository{DbContext: dbContext}
}

func (repo ProductRepository) GetProducts() ([]domain.Product, error) {
	var products []domain.Product
	if err := repo.DbContext.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repo ProductRepository) GetProductBySku(sku string) (*domain.Product, error) {

	prd := domain.Product{}
	if err := repo.DbContext.Where(&domain.Product{Sku: sku}).First(&prd).Error; err != nil {
		return nil, err
	}
	return &prd, nil
}

func (repo ProductRepository) SaveStock(country, sku, name string, stockChange int) (bool, error) {
	var prd *domain.Product = &domain.Product{}

	if err := repo.DbContext.Where(&domain.Product{Sku: sku}).First(&prd).Error; err == nil {
		if ok := math.Signbit(float64(stockChange)); ok {
			absValue := int(math.Abs(float64(stockChange)))
			if _, err := prd.RemoveStock(absValue); err != nil {
				return false, err
			}
		} else {
			if _, err := prd.AddStock(stockChange); err != nil {
				return false, err
			}
		}
		if err := repo.DbContext.Save(&prd).Error; err != nil {
			return false, err
		}
	} else {
		newProduct := domain.Product{
			Name:           name,
			Sku:            sku,
			Country:        country,
			AvailableStock: 0,
		}

		if err := repo.DbContext.Create(&newProduct).Error; err != nil {
			return false, err
		}
	}

	return true, nil
}
