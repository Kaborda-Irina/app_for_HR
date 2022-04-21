package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/inkoba/app_for_HR/internal/config"
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/domain/request"
	"github.com/inkoba/app_for_HR/internal/core/domain/response"
	mock_ports "github.com/inkoba/app_for_HR/internal/core/ports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSalaryService_Create(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockISalaryRepository, file []byte, salaries []*domain.Salary)

	testTable := []struct {
		name          string
		file          []byte
		mockBehavior  mockBehavior
		inputData     []*domain.Salary
		expected      *response.SalaryUploadReport
		expectedError bool
	}{
		{
			name: "get error when delete all data in database is impossible",
			file: []byte{},
			mockBehavior: func(s *mock_ports.MockISalaryRepository, file []byte, salaries []*domain.Salary) {
				s.EXPECT().DeleteAll().Return(errors.New("Error delete all data in database"))
			},
			expectedError: true,
		},
		{
			name: "get error when creating new salaries in database and database is unavailable",
			file: []byte{},
			mockBehavior: func(s *mock_ports.MockISalaryRepository, file []byte, salaries []*domain.Salary) {
				s.EXPECT().DeleteAll()
				s.EXPECT().Create(salaries).Return(errors.New("Error delete all data in database"))
			},
			expectedError: true,
		},
		{
			name: "creating report is successful",
			file: []byte{},
			mockBehavior: func(s *mock_ports.MockISalaryRepository, file []byte, salaries []*domain.Salary) {
				s.EXPECT().DeleteAll()
				s.EXPECT().Create(salaries)
			},
			expected: &response.SalaryUploadReport{
				0,
				0,
				[]string(nil),
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_ports.NewMockISalaryRepository(c)
			testCase.mockBehavior(repo, testCase.file, testCase.inputData)
			newConfig := config.CurrencyConfig{CoefficientEURtoUSD: 1.1328, CoefficientRUStoUSD: 0.014}
			service := SalaryService{repo, newConfig, logrus.New()}

			wantResult, err := service.Create(testCase.file)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, wantResult)
			}
		})
	}
}

func TestSalaryService_GetSalariesByFilter(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockISalaryRepository, salaryFilteringCondition *request.ConditionForFilteringSalaries)

	testTable := []struct {
		name                     string
		salaryFilteringCondition *request.ConditionForFilteringSalaries
		mockBehavior             mockBehavior
		expected                 []*response.SalariesResponse
		expectedError            bool
	}{
		{
			name: "salary filter works successfully when all field is filled",
			salaryFilteringCondition: &request.ConditionForFilteringSalaries{
				Salary:           "1500",
				LevelOfSeniority: "Junior",
				YearsTotal:       "1",
				Country:          "Belarus",
			},
			mockBehavior: func(s *mock_ports.MockISalaryRepository, salaryFilteringCondition *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetFilteredSalaries(salaryFilteringCondition).Return([]*domain.Salary{}, nil)
			},
			expected: []*response.SalariesResponse(nil),
		},

		{
			name: "salary filter can not filtering when database is unavailable",
			salaryFilteringCondition: &request.ConditionForFilteringSalaries{
				Salary:           "1500",
				LevelOfSeniority: "Junior",
				YearsTotal:       "1",
				Country:          "Belarus",
			},
			mockBehavior: func(s *mock_ports.MockISalaryRepository, salaryFilteringCondition *request.ConditionForFilteringSalaries) {
				s.EXPECT().GetFilteredSalaries(salaryFilteringCondition).Return(nil, errors.New("Error get filtered salaries "))
			},
			expectedError: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_ports.NewMockISalaryRepository(c)
			testCase.mockBehavior(repo, testCase.salaryFilteringCondition)
			newConfig := config.CurrencyConfig{CoefficientEURtoUSD: 1.1328, CoefficientRUStoUSD: 0.014}
			service := SalaryService{repo, newConfig, logrus.New()}

			wantResult, err := service.GetSalariesByFilter(testCase.salaryFilteringCondition)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, wantResult)
			}
		})
	}
}
