package models

// ResponseModelList :
type ResponseModelList struct {
	Page     int         `json:"page"`
	Total    int64       `json:"total"`
	LastPage int64       `json:"last_page"`
	Status   int         `json:"status"`
	Msg      string      `json:"message"`
	Data     interface{} `json:"data"`
}
