package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	mock_ports "github.com/inkoba/app_for_HR/internal/core/ports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthService_Ping(t *testing.T) {
	testTable := []struct {
		name          string
		mockBehavior  func(s *mock_ports.MockIHealthService)
		expected      interface{}
		expectedError bool
	}{
		{
			name: "database connected successfully",
			mockBehavior: func(s *mock_ports.MockIHealthService) {
				s.EXPECT().Ping().Return(nil)
			},
			expected: nil,
		},
		{
			name: "database is not connected",
			mockBehavior: func(s *mock_ports.MockIHealthService) {
				s.EXPECT().Ping().Return(errors.New("database is not connected"))
			},
			expectedError: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_ports.NewMockIHealthService(c)
			testCase.mockBehavior(repo)

			service := HealthService{repo, logrus.New()}

			err := service.Ping()

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
