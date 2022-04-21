package request

type ConditionForFilteringSalaries struct {
	Salary           string `json:"salary"`
	LevelOfSeniority string `json:"levelOfSeniority"`
	YearsTotal       string `json:"yearsTotal"`
	Country          string `json:"country"`
}
