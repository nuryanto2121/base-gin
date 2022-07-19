package models

// ResponseModelList :
type ResponseModelList struct {
	Page         int         `json:"page"`
	Total        int64       `json:"total"`
	LastPage     int64       `json:"last_page"`
	DefineSize   string      `json:"define_size"`
	DefineColumn string      `json:"define_column"`
	AllColumn    string      `json:"all_column"`
	Version      string      `json:"version"`
	Status       int         `json:"status"`
	Msg          string      `json:"message"`
	Data         interface{} `json:"data"`
}
