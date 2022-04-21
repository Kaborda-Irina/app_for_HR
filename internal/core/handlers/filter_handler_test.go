package handlers

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/inkoba/app_for_HR/internal/core/domain/request"
	"github.com/inkoba/app_for_HR/internal/core/domain/response"
	mock_ports "github.com/inkoba/app_for_HR/internal/core/ports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestSalaryFilterHandler_Filter(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries)
	testTable := []struct {
		name                 string
		inputBody            string
		inputCondition       request.ConditionForFilteringSalaries
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "salary filter works successfully when one field is filled",
			inputBody: `{"country":"","salary":"1500","yearsTotal":"","levelOfSeniority":""}`,
			inputCondition: request.ConditionForFilteringSalaries{
				Salary:           "1500",
				LevelOfSeniority: "",
				YearsTotal:       "",
				Country:          ""},
			mockBehavior: func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetSalariesByFilter(salaries).
					Return([]*response.SalariesResponse{{
						Salary:           "1500",
						LevelOfSeniority: "Junior",
						YearsTotal:       "1",
						Country:          "Belarus"},
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"}]
`,
		},
		{
			name:      "salary filter works successfully when two fields are filled",
			inputBody: `{"country":"","salary":"1500","yearsTotal":"","levelOfSeniority":"Junior"}`,
			inputCondition: request.ConditionForFilteringSalaries{
				Salary:           "1500",
				LevelOfSeniority: "Junior",
				YearsTotal:       "",
				Country:          ""},
			mockBehavior: func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetSalariesByFilter(salaries).
					Return([]*response.SalariesResponse{{
						Salary:           "1500",
						LevelOfSeniority: "Junior",
						YearsTotal:       "1",
						Country:          "Belarus"}}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"}]
`,
		},
		{
			name:      "salary filter works successfully when three fields are filled",
			inputBody: `{"country":"Belarus","salary":"1500","yearsTotal":"","levelOfSeniority":"Junior"}`,
			inputCondition: request.ConditionForFilteringSalaries{
				Salary:           "1500",
				LevelOfSeniority: "Junior",
				YearsTotal:       "",
				Country:          "Belarus"},
			mockBehavior: func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetSalariesByFilter(salaries).
					Return([]*response.SalariesResponse{{
						Salary:           "1500",
						LevelOfSeniority: "Junior",
						YearsTotal:       "1",
						Country:          "Belarus"}}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"}]
`,
		},
		{
			name:      "salary filter works successfully when all fields are filled",
			inputBody: `{"country":"Belarus","salary":"1500","yearsTotal":"1","levelOfSeniority":"Junior"}`,
			inputCondition: request.ConditionForFilteringSalaries{
				Salary:           "1500",
				LevelOfSeniority: "Junior",
				YearsTotal:       "1",
				Country:          "Belarus"},
			mockBehavior: func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetSalariesByFilter(salaries).
					Return([]*response.SalariesResponse{{
						Salary:           "1500",
						LevelOfSeniority: "Junior",
						YearsTotal:       "1",
						Country:          "Belarus"},
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"}]
`,
		},
		{
			name:      "salary filter works successfully when all fields are empty",
			inputBody: `{"country":"","salary":"","yearsTotal":"","levelOfSeniority":""}`,
			inputCondition: request.ConditionForFilteringSalaries{
				Salary:           "",
				LevelOfSeniority: "",
				YearsTotal:       "",
				Country:          ""},
			mockBehavior: func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetSalariesByFilter(salaries).
					Return([]*response.SalariesResponse{{
						Salary:           "1500",
						LevelOfSeniority: "Junior",
						YearsTotal:       "1",
						Country:          "Belarus",
					},
						{
							Salary:           "1500",
							LevelOfSeniority: "Junior",
							YearsTotal:       "1",
							Country:          "Belarus"},
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"},{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"}]
`,
		},
		{
			name:      "salary filter works successfully when fields contain spaces",
			inputBody: `{"country":"   ","salary":"  ","yearsTotal":"","levelOfSeniority":""}`,
			inputCondition: request.ConditionForFilteringSalaries{
				Salary:           "  ",
				LevelOfSeniority: "",
				YearsTotal:       "",
				Country:          "   "},
			mockBehavior: func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetSalariesByFilter(salaries).
					Return([]*response.SalariesResponse{{
						Salary:           "1500",
						LevelOfSeniority: "Junior",
						YearsTotal:       "1",
						Country:          "Belarus",
					},
						{
							Salary:           "1500",
							LevelOfSeniority: "Junior",
							YearsTotal:       "1",
							Country:          "Belarus"},
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"},{"salary":"1500","levelOfSeniority":"Junior","yearsTotal":"1","country":"Belarus"}]
`,
		},
		{
			name:               "salary filter has no string value in the country field then an error follows",
			inputBody:          `{"country":1213,"salary":"","yearsTotal":"","levelOfSeniority":""}`,
			inputCondition:     request.ConditionForFilteringSalaries{},
			mockBehavior:       func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["json: cannot unmarshal number into Go struct field ConditionForFilteringSalaries.country of type string"]}
`,
		},
		{
			name:      "salary filter can't filtered data in database when the database fail",
			inputBody: `{"country":"Belarus","salary":"","yearsTotal":"","levelOfSeniority":""}`,
			inputCondition: request.ConditionForFilteringSalaries{
				Salary:           "",
				LevelOfSeniority: "",
				YearsTotal:       "",
				Country:          "Belarus",
			},
			mockBehavior: func(s *mock_ports.MockISalaryService, salaries *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetSalariesByFilter(salaries).
					Return(nil, errors.New("Error getting filtered salaries"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"Errors":["Error getting filtered salaries"]}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockISalaryService(c)
			testCase.mockBehavior(service, &testCase.inputCondition)

			handler := SalaryFilterHandler{service, logrus.New()}

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/api/filter", handler.Filter).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/filter",
				bytes.NewBufferString(testCase.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
