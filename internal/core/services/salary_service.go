package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/inkoba/app_for_HR/internal/config"
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/domain/request"
	"github.com/inkoba/app_for_HR/internal/core/domain/response"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
)

const (
	currencyUSD = "USD"
	currencyEUR = "EUR"
	currencyRUS = "RUS"

	bitSize = 64
)
const (
	twoElementsSalary   = 2
	firstElementSalary  = 0
	secondElementSalary = 1
)

const (
	indexSalary           = 23
	indexLevelOfSeniority = 27
	indexYearsTotal       = 29
	indexCountry          = 37
	indexLevelOfEnglish   = 25
)

type SalaryService struct {
	salaryRepository ports.ISalaryRepository
	currencyConfig   config.CurrencyConfig
	logger           *logrus.Logger
}

var _ ports.ISalaryService = (*SalaryService)(nil)

func NewSalaryService(currencyConfig config.CurrencyConfig, repository ports.ISalaryRepository, logger *logrus.Logger) *SalaryService {
	return &SalaryService{
		repository,
		currencyConfig,
		logger,
	}
}

func (ss SalaryService) Create(file []byte) (*response.SalaryUploadReport, error) {

	//Clean up salaries from storage
	err := ss.salaryRepository.DeleteAll()
	if err != nil {
		ss.logger.Error(err)
		return nil, err
	}

	ss.logger.Info("Old documents in collection deleted successfully")

	lines, err := csv.NewReader(bytes.NewReader(file)).ReadAll()
	if err != nil {
		ss.logger.Error(err)
		return nil, err
	}
	report := response.SalaryUploadReport{}
	var salaries []*domain.Salary
	for index, line := range lines {
		if index == 0 {
			continue
		}

		report.TotalRecords++

		if isLineNotValid(line) {
			errorHandler(&report, fmt.Errorf("Line %d is not valid", index))
			continue
		}

		salaryElements := strings.Split(line[indexSalary], " ")

		if len(salaryElements) == twoElementsSalary {
			for index, elem := range salaryElements {
				if index == firstElementSalary {
					salaryElements[firstElementSalary] = elem
				} else {
					salaryElements[secondElementSalary] = elem
				}
			}
		} else {
			err := fmt.Errorf("The salary field has more elements than being processed")
			errorHandler(&report, err)
			continue
		}

		line[indexSalary] = strings.Join(salaryElements, " ")

		userDataOnSalary := domain.Salary{
			Salary:           salaryElements[firstElementSalary],
			Currency:         salaryElements[secondElementSalary],
			LevelOfSeniority: line[indexLevelOfSeniority],
			YearsTotal:       line[indexYearsTotal],
			Country:          line[indexCountry],
			LevelOfEnglish:   line[indexLevelOfEnglish],
		}
		salaries = append(salaries, &userDataOnSalary)
	}

	err = ss.salaryRepository.Create(salaries)
	return &report, err
}

func (ss SalaryService) GetSalariesByFilter(salaryFilteringCondition *request.ConditionForFilteringSalaries) ([]*response.SalariesResponse, error) {
	filteredSalaries, err := ss.salaryRepository.GetFilteredSalaries(salaryFilteringCondition)
	if err != nil {
		ss.logger.Error("Error get filtered salaries: ")
		return nil, err
	}
	var filteredResponse []*response.SalariesResponse
	for _, filteredElem := range filteredSalaries {

		if isCurrensyNotValid(filteredElem, ss) {
			continue
		}

		result := response.SalariesResponse{
			Salary:           filteredElem.Salary,
			LevelOfSeniority: filteredElem.LevelOfSeniority,
			YearsTotal:       filteredElem.YearsTotal,
			Country:          filteredElem.Country,
		}

		filteredResponse = append(filteredResponse, &result)
	}
	return filteredResponse, nil
}

func errorHandler(report *response.SalaryUploadReport, err error) {
	report.Errors = append(report.Errors, err.Error())
	report.SkippedRecords++
}

func isLineNotValid(line []string) bool {
	return isIndexOfSalaryEmpty(line) ||
		len(strings.TrimSpace(line[indexLevelOfSeniority])) == 0 ||
		len(strings.TrimSpace(line[indexYearsTotal])) == 0 ||
		len(strings.TrimSpace(line[indexCountry])) == 0 ||
		len(strings.TrimSpace(line[indexLevelOfEnglish])) == 0
}

func isIndexOfSalaryEmpty(line []string) bool {
	return len(strings.TrimSpace(line[indexSalary])) == 0
}
func isCurrensyNotValid(filteredElem *domain.Salary, ss SalaryService) bool {
	return isCurrencyNotConvertible(filteredElem, ss)
}
func isCurrencyNotConvertible(filteredElem *domain.Salary, ss SalaryService) bool {
	switch filteredElem.Currency {
	case currencyUSD:
	case currencyEUR:
		converter(filteredElem, ss.logger, ss.currencyConfig.CoefficientEURtoUSD)
	case currencyRUS:
		converter(filteredElem, ss.logger, ss.currencyConfig.CoefficientRUStoUSD)
	default:
		ss.logger.Error("CurrencyConfig does not exist for line: ", filteredElem)
		return true
	}
	return false
}

func converter(salary *domain.Salary, logger *logrus.Logger, coefficient float64) {
	currency, err := strconv.ParseFloat(salary.Salary, bitSize)
	if err != nil {
		logger.Error(err)
	}
	salary.Salary = fmt.Sprint(math.Round(currency * coefficient))
}
