package entity

type Result struct {
	ID        string `json:"id"`
	TaskKey   string `json:"taskKey"`
	Total     string `json:"total"`
	Tenant    string `json:"tenant"`
	StartUser string `json:"startUser"`
	StartTime string `json:"startTime"`
	BussKey   string `json:"bussKey"`
}
