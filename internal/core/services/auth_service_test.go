package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/inkoba/app_for_HR/internal/core/domain"
	mock_ports "github.com/inkoba/app_for_HR/internal/core/ports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestAuthService_IsValidUser(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, username string, password string)
	testTable := []struct {
		name          string
		username      string
		password      string
		mockBehavior  mockBehavior
		expected      *domain.User
		expectedError bool
	}{
		{
			name:     "username exists in database",
			username: "admin",
			password: "1234",
			mockBehavior: func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, username string, password string) {
				s.EXPECT().GetUserByUsername(username).Return(&domain.User{
					Id:       primitive.ObjectID{},
					Username: "admin",
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  false,
				}, nil)
				h.EXPECT().CompareHashAndPassword([]byte("$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq"), []byte(password)).Return(nil)
			},
			expected: nil,
		},
		{
			name:     "username does not exist in database",
			username: "vasya",
			password: "1234",
			mockBehavior: func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, username string, password string) {
				s.EXPECT().GetUserByUsername(username).Return(nil, errors.New(": username does not exist in the database"))
			},
			expectedError: true,
		},
		{
			name:     "password in the database does not match hash password",
			username: "admin",
			password: "1234",
			mockBehavior: func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, username string, password string) {
				s.EXPECT().GetUserByUsername(username).Return(&domain.User{
					Id:       primitive.ObjectID{},
					Username: "admin",
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  false,
				}, nil)
				h.EXPECT().CompareHashAndPassword([]byte("$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq"), []byte(password)).Return(errors.New("password in the database does not match hash password"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_ports.NewMockIUserRepository(c)
			repoCrypto := mock_ports.NewMockICryptoService(c)
			testCase.mockBehavior(repo, repoCrypto, testCase.username, testCase.password)

			service := AuthService{logrus.New(), repo, repoCrypto}

			err := service.IsValidUser(testCase.username, testCase.password)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
