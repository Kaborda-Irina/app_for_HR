package response

type SalaryUploadReport struct {
	TotalRecords   int
	SkippedRecords int
	Errors         []string
}
