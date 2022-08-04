package repositories

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/philaden/mds-stock-keeping/application/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetProducts(t *testing.T) {

	products := mocks.GetMockProducts()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepository := mocks.NewMockIProductRepository(ctrl)
	mockProductRepository.EXPECT().GetProducts().Return(products, nil)

	response, err := mockProductRepository.GetProducts()
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, response, products)
}

func TestGetProductBySku(t *testing.T) {
	const sku string = "da8ef851e075"
	product := mocks.GetMockProductBySku(sku)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepository := mocks.NewMockIProductRepository(ctrl)
	mockProductRepository.EXPECT().GetProductBySku(sku).Return(product, nil)

	response, err := mockProductRepository.GetProductBySku(sku)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, response, product)
}

func TestCreateStock(t *testing.T) {
	data := mocks.CreateMockStockPayload()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepository := mocks.NewMockIProductRepository(ctrl)
	mockProductRepository.EXPECT().SaveStock(data.Country, data.Sku, data.Name, data.StockChange).Return(true, nil)

	response, err := mockProductRepository.SaveStock(data.Country, data.Sku, data.Name, data.StockChange)
	require.NoError(t, err)
	require.Equal(t, response, true)
}
