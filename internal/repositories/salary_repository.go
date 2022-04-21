package repositories

import (
	"context"
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/domain/request"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
)

type SalaryRepository struct {
	mc     *MongoConfig
	logger *logrus.Logger
}

var _ ports.ISalaryRepository = (*SalaryRepository)(nil)

func NewSalaryRepository(mc *MongoConfig, logger *logrus.Logger) ports.ISalaryRepository {
	return &SalaryRepository{
		mc,
		logger,
	}
}

func (sr SalaryRepository) Create(salaries []*domain.Salary) error {
	result := bson.A{}
	for _, data := range salaries {
		result = append(result, data)
	}

	_, err := sr.mc.salariesCollection.InsertMany(context.Background(), result)
	if err != nil {
		return err
	}

	return err
}

func (sr SalaryRepository) DeleteAll() error {
	return sr.mc.salariesCollection.Drop(context.Background())
}

func (sr SalaryRepository) GetFilteredSalaries(salaryFilteringCondition *request.ConditionForFilteringSalaries) ([]*domain.Salary, error) {
	filter := filteredFields(salaryFilteringCondition)
	cursor, err := sr.mc.salariesCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			sr.logger.Error(err)
		}
	}(cursor, context.Background())

	var list []*domain.Salary
	for cursor.Next(context.Background()) {
		var salary *domain.Salary

		err := cursor.Decode(&salary)
		if err != nil {
			return nil, err
		}

		result := &domain.Salary{
			Salary:           salary.Salary,
			LevelOfSeniority: salary.LevelOfSeniority,
			YearsTotal:       salary.YearsTotal,
			Country:          salary.Country,
			Currency:         salary.Currency,
		}

		list = append(list, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func filteredFields(filterSalary *request.ConditionForFilteringSalaries) bson.M {
	filter := bson.M{}

	if len(strings.TrimSpace(filterSalary.Salary)) > 0 {
		filter["salary"] = filterSalary.Salary
	}
	if len(strings.TrimSpace(filterSalary.Country)) > 0 {
		filter["country"] = filterSalary.Country
	}
	if len(strings.TrimSpace(filterSalary.YearsTotal)) > 0 {
		filter["yearstotal"] = filterSalary.YearsTotal
	}
	if len(strings.TrimSpace(filterSalary.LevelOfSeniority)) > 0 {
		filter["levelofseniority"] = filterSalary.LevelOfSeniority
	}

	return filter
}
