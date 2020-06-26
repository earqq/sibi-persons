package model

type Purchase struct {
	Serie       string  `json:"serie"`
	Number      int64   `json:"number"`
	CompanyRuc  string  `json:"company_ruc" bson:"company_ruc"`
	CompanyName string  `json:"company_name" bson:"company_name"`
	TotalPrice  float64 `json:"total_price" bson:"total_price"`
	TotalIgv    float64 `json:"total_igv" bson:"total_igv"`
	IssueDate   string  `json:"issue_date" bson:"issue_date"`
}
