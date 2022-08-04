package repositories

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/philaden/mds-stock-keeping/application/mocks"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	data := mocks.CreateMockOrderPayload()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepository := mocks.NewMockIOrderRepository(ctrl)
	mockOrderRepository.EXPECT().CreateSingleOrder(data).Return(data.ID, nil)

	response, err := mockOrderRepository.CreateSingleOrder(data)
	require.NoError(t, err)
	require.Equal(t, response, data.ID)
}
