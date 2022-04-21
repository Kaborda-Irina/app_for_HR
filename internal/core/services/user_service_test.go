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

var id = [12]byte{61, 'b', 73, 4, 137, 8, 'a', 'd', 60, 'a', 0, 'd'}

func TestUserService_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserRepository)
	testTable := []struct {
		name          string
		mockBehavior  mockBehavior
		expected      []*domain.User
		expectedError bool
	}{
		{
			name: "get all users when the database is available",
			mockBehavior: func(s *mock_ports.MockIUserRepository) {
				s.EXPECT().GetAll().Return([]*domain.User{
					{
						Id:       id,
						Username: "admin",
						Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
						IsAdmin:  true,
					},
					{
						Id:       id,
						Username: "user",
						Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
						IsAdmin:  false,
					},
				}, nil)
			},
			expected: []*domain.User{
				{
					Id:       id,
					Username: "admin",
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  true,
				},
				{
					Id:       id,
					Username: "user",
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  false,
				},
			},
		},
		{
			name: "get all users when the database is unavailable",
			mockBehavior: func(s *mock_ports.MockIUserRepository) {
				s.EXPECT().GetAll().Return(nil, errors.New("database is unavailable"))
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
			testCase.mockBehavior(repo)

			service := UserService{repo, logrus.New(), repoCrypto}

			wantResult, err := service.GetAll()

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, wantResult)
			}

		})
	}
}

func TestUserService_Get(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserRepository, idUser string)
	testTable := []struct {
		name          string
		idUser        string
		mockBehavior  mockBehavior
		expected      *domain.User
		expectedError bool
	}{
		{
			name:   "get one user from the database when  the user id is correct",
			idUser: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserRepository, idUser string) {
				s.EXPECT().Get(idUser).Return(&domain.User{
					Id:       id,
					Username: "admin",
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  true,
				}, nil)
			},
			expected: &domain.User{
				Id:       id,
				Username: "admin",
				Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
				IsAdmin:  true,
			}},
		{
			name:   "get error when the user id is incorrect",
			idUser: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserRepository, idUser string) {
				s.EXPECT().Get(idUser).Return(nil, errors.New(`"Errors": ["mongo: no documents in result"]`))
			},
			expectedError: true},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_ports.NewMockIUserRepository(c)
			repoCrypto := mock_ports.NewMockICryptoService(c)
			testCase.mockBehavior(repo, testCase.idUser)

			service := UserService{repo, logrus.New(), repoCrypto}

			wantResult, err := service.Get(testCase.idUser)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, wantResult)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserRepository, idUser string)
	testTable := []struct {
		name          string
		idUser        string
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name:   "delete one user is successful",
			idUser: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserRepository, idUser string) {
				s.EXPECT().Delete(idUser)
			},
		},
		{
			name:   "delete one user when user data not in database",
			idUser: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserRepository, idUser string) {
				s.EXPECT().Delete(idUser).Return(errors.New(`"Errors": ["there is no such user in the database"]`))
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
			testCase.mockBehavior(repo, testCase.idUser)

			service := UserService{repo, logrus.New(), repoCrypto}

			err := service.Delete(testCase.idUser)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestUserService_Create(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, user *domain.User)

	testTable := []struct {
		name          string
		inputData     *domain.User
		mockBehavior  mockBehavior
		expected      string
		expectedError bool
	}{
		{
			name: "create user with data which not exist in database",
			inputData: &domain.User{
				Id:       primitive.ObjectID{},
				IsAdmin:  false,
				Password: "1111",
				Username: "admin",
			},
			mockBehavior: func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, user *domain.User) {
				h.EXPECT().GetHashedPassword([]byte(user.Password)).Return("$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq", nil)
				newUser := &domain.User{
					Username: user.Username,
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  false,
				}
				s.EXPECT().Create(newUser).Return("000000000000000000000000", nil)
			},
			expected: "000000000000000000000000",
		},
		{
			name: "error when the password can't be hashed",
			inputData: &domain.User{
				Id:       primitive.ObjectID{},
				IsAdmin:  false,
				Password: "1111",
				Username: "admin",
			},
			mockBehavior: func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, user *domain.User) {
				h.EXPECT().GetHashedPassword([]byte(user.Password)).Return("", errors.New("password cannot be hashed"))
			},
			expectedError: true,
		},
		{
			name: "create user with data which not exist in database",
			inputData: &domain.User{
				Id:       primitive.ObjectID{},
				IsAdmin:  false,
				Password: "1111",
				Username: "admin",
			},
			mockBehavior: func(s *mock_ports.MockIUserRepository, h *mock_ports.MockICryptoService, user *domain.User) {
				h.EXPECT().GetHashedPassword([]byte(user.Password)).Return("$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq", nil)
				newUser := &domain.User{
					Username: user.Username,
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  false,
				}
				s.EXPECT().Create(newUser).Return("", errors.New("new user can not be created"))
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
			testCase.mockBehavior(repo, repoCrypto, testCase.inputData)

			service := UserService{repo, logrus.New(), repoCrypto}

			wantResult, err := service.Create(testCase.inputData)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, wantResult)
			}
		})
	}
}

func TestUserService_GetUserByUsername(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserRepository, username string)

	testTable := []struct {
		name          string
		inputData     *domain.User
		username      string
		mockBehavior  mockBehavior
		expected      *domain.User
		expectedError bool
	}{
		{
			name: "get user when user exist in database",
			inputData: &domain.User{
				Id:       id,
				IsAdmin:  false,
				Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
				Username: "admin",
			},
			username: "admin",
			mockBehavior: func(s *mock_ports.MockIUserRepository, username string) {
				s.EXPECT().GetUserByUsername(username).Return(
					&domain.User{
						Id:       id,
						IsAdmin:  false,
						Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
						Username: "admin",
					}, nil)
			},
			expected: &domain.User{
				Id:       id,
				IsAdmin:  false,
				Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
				Username: "admin",
			},
		},

		{
			name:     "get user when user does not exist in database",
			username: "vasya/45",
			mockBehavior: func(s *mock_ports.MockIUserRepository, username string) {
				s.EXPECT().GetUserByUsername(username).Return(nil, errors.New("Errors: user is not exist in database"))
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
			testCase.mockBehavior(repo, testCase.username)

			service := UserService{repo, logrus.New(), repoCrypto}

			wantResult, err := service.GetUserByUsername(testCase.username)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, wantResult)
			}
		})
	}
}
