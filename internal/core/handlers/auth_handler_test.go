package handlers

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	mock_ports "github.com/inkoba/app_for_HR/internal/core/ports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_Login_ValidUser(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIAuthService, username string, password string)
	testTable := []struct {
		name                 string
		inputBody            string
		inputCredentials     Credentials
		username             string
		password             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:               "unable to decode request body",
			username:           "user",
			password:           "1234",
			mockBehavior:       func(s *mock_ports.MockIAuthService, username string, password string) {},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["EOF"]}
`,
		},
		{
			name:      "get error when not valid data in request",
			inputBody: `{"username":"user/4844","password":"1234/56219**"}`,
			inputCredentials: Credentials{
				Username: "user/4844",
				Password: "1234/56219**",
			},
			username: "user/4844",
			password: "1234/56219**",
			mockBehavior: func(s *mock_ports.MockIAuthService, username string, password string) {
				s.EXPECT().IsValidUser(username, password).Return(errors.New("Error validation username or password"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["Error validation username or password"]}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			serviceAuth := mock_ports.NewMockIAuthService(c)
			testCase.mockBehavior(serviceAuth, testCase.username, testCase.password)

			serviceUser := mock_ports.NewMockIUserService(c)

			handler := AuthHandler{serviceAuth, serviceUser, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/login", handler.Login).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/login",
				bytes.NewBufferString(testCase.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestAuthHandler_Login_GetUserByUsername(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIUserService, sa *mock_ports.MockIAuthService, username string, password string)
	testTable := []struct {
		name                 string
		inputBody            string
		username             string
		password             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "username does not exit in database",
			inputBody: `{"username":"admin","password":"1234"}`,
			username:  "admin",
			password:  "1234",
			mockBehavior: func(s *mock_ports.MockIUserService, sa *mock_ports.MockIAuthService, username string, password string) {
				sa.EXPECT().IsValidUser(username, password).Return(nil)
				s.EXPECT().GetUserByUsername(username).Return(nil, errors.New("User does not exist in database"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["User does not exist in database"]}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			serviceAuth := mock_ports.NewMockIAuthService(c)
			serviceUser := mock_ports.NewMockIUserService(c)
			testCase.mockBehavior(serviceUser, serviceAuth, testCase.username, testCase.password)
			handler := AuthHandler{serviceAuth, serviceUser, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/login", handler.Login).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/login",
				bytes.NewBufferString(testCase.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
