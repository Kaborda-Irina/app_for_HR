package handlers

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/inkoba/app_for_HR/internal/core/domain"
	mock_ports "github.com/inkoba/app_for_HR/internal/core/ports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var id = [12]byte{61, 'b', 73, 4, 137, 8, 'a', 'd', 60, 'a', 0, 'd'}

func TestUserHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserService)
	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "get all users when the database is available",
			mockBehavior: func(s *mock_ports.MockIUserService) {
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
			expectedStatusCode: 200,
			expectedResponseBody: `[{"id":"3d624904890861643c610064","username":"admin","password":"$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq","isAdmin":true},{"id":"3d624904890861643c610064","username":"user","password":"$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq","isAdmin":false}]
`},
		{
			name: "get error when the database is unavailable",
			mockBehavior: func(s *mock_ports.MockIUserService) {
				s.EXPECT().GetAll().Return(nil, errors.New("database is unavailable"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["database is unavailable"]}
null
`},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockIUserService(c)
			testCase.mockBehavior(service)

			handler := UserHandler{service, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/users", handler.GetAll).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/users", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserService, inputId string)

	testTable := []struct {
		name                 string
		inputId              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "get one user from the database when user id is correct",
			inputId: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserService, inputId string) {
				s.EXPECT().Get(inputId).Return(&domain.User{
					Id:       id,
					Username: "admin",
					Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
					IsAdmin:  true,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"id":"3d624904890861643c610064","username":"admin","password":"$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq","isAdmin":true}
`},
		{
			name:    "get error when the user id is incorrect",
			inputId: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserService, id string) {
				s.EXPECT().Get(id).Return(nil, errors.New(`"Errors": ["mongo: no documents in result"]`))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["\"Errors\": [\"mongo: no documents in result\"]"]}
null
`},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockIUserService(c)
			testCase.mockBehavior(service, testCase.inputId)

			handler := UserHandler{service, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/users/{id:[a-zA-Z0-9]*}", handler.Get).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/users/3d624904890861643c610064", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestUserHandler_Create_GetUserByUsername(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserService, username string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            *domain.User
		username             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "new user is exist in database",
			inputBody: `{"id":"3d624904890861643c610064","isAdmin":false,"password": "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq","username":"admin"}`,
			inputData: &domain.User{
				Id:       id,
				IsAdmin:  false,
				Password: "$2a$12$QSEvrvXWWegdupNz73bYeedLkOl5VRUNWT8iG2hGeeN5Z1FjlfBxq",
				Username: "admin",
			},
			username: "admin",
			mockBehavior: func(s *mock_ports.MockIUserService, username string) {
				s.EXPECT().GetUserByUsername(username).Return(&domain.User{}, nil)
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["User is exist in database"]}
`,
		},

		{
			name:               "get error when unable to decode request body",
			username:           "vasya/45",
			mockBehavior:       func(s *mock_ports.MockIUserService, username string) {},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["EOF"]}
`},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockIUserService(c)
			testCase.mockBehavior(service, testCase.username)

			handler := UserHandler{service, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/users", handler.Create).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(testCase.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestUserHandler_Create_CreateNewUser(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserService, user *domain.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            *domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "create user with data which not exist in database",
			inputBody: `{"id":"3d624904890861643c610064","isAdmin":false,"password": "1234","username":"admin"}`,
			inputData: &domain.User{
				Id:       id,
				IsAdmin:  false,
				Password: "1234",
				Username: "admin",
			},
			mockBehavior: func(s *mock_ports.MockIUserService, user *domain.User) {
				s.EXPECT().GetUserByUsername(user.Username).Return(nil, errors.New("User does not exist in database"))
				s.EXPECT().Create(user).Return("3d624904890861643c610064", nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `"3d624904890861643c610064"
`},
		{
			name:      "get an error when creating a user with data that exists in the database",
			inputBody: `{"id":"3d624904890861643c610064","isAdmin":false,"password": "1234","username":"admin"}`,
			inputData: &domain.User{
				Id:       id,
				IsAdmin:  false,
				Password: "1234",
				Username: "admin",
			},
			mockBehavior: func(s *mock_ports.MockIUserService, user *domain.User) {
				s.EXPECT().GetUserByUsername(user.Username).Return(nil, errors.New("User does not exist in database"))
				s.EXPECT().Create(user).Return("", errors.New("error create new user in database"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["error create new user in database"]}
""
`},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockIUserService(c)
			testCase.mockBehavior(service, testCase.inputData)

			handler := UserHandler{service, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/users", handler.Create).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(testCase.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserService, inputId string)

	testTable := []struct {
		name                 string
		inputId              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "delete one user from the database with valid data",
			inputId: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserService, inputId string) {
				s.EXPECT().Delete(inputId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:    "delete one user when user data not in database",
			inputId: "3d624904890861643c610064",
			mockBehavior: func(s *mock_ports.MockIUserService, inputId string) {
				s.EXPECT().Delete(inputId).Return(errors.New(`"Errors": ["there is no such user in the database"]`))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["\"Errors\": [\"there is no such user in the database\"]"]}
`},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockIUserService(c)
			testCase.mockBehavior(service, testCase.inputId)

			handler := UserHandler{service, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/users/{id:[a-zA-Z0-9]*}", handler.Delete).Methods("DELETE")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/users/3d624904890861643c610064", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
