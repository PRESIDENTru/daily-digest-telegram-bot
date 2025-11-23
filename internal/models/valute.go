package models

type Valute struct {
	Rates map[string]float64 `json:"rates"`
	Err   error
	Code  string
}
