package model

type Purchase struct {
	Serie           string  `json:"serie"`
	Number          int64   `json:"number"`
	ContactIdentity string  `json:"contact_identity" bson:"contact_identity`
	ContactName     string  `json:"contact_name" bson:"contact_name"`
	TotalPrice      float64 `json:"total_price" bson:"total_price"`
	TotalIgv        float64 `json:"total_igv" bson:"total_igv"`
	IssueDate       string  `json:"issue_date" bson:"issue_date"`
}
