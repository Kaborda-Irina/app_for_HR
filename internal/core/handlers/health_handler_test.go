package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	mock_ports "github.com/inkoba/app_for_HR/internal/core/ports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_Ping(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIHealthService)
	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "health check when there is a connection to the mongodb database",
			mockBehavior: func(s *mock_ports.MockIHealthService) {
				s.EXPECT().Ping().Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `API is up and running`,
		},
		{
			name: "health check when there is no connection to mongodb database",
			mockBehavior: func(s *mock_ports.MockIHealthService) {
				s.EXPECT().Ping().Return(errors.New("Mongo is NOT up and running"))
			},
			expectedStatusCode:   200,
			expectedResponseBody: `Mongo is NOT up and running`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockIHealthService(c)
			testCase.mockBehavior(service)

			handler := HealthHandler{service, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/health", handler.Ping).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/health", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
