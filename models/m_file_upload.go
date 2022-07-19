package models

type FileUpload struct {
	ID       int    `json:"id" gorm:"primary_key;auto_increment:true"`
	FileName string `json:"file_name" gorm:"type:varchar(60)"`
	FilePath string `json:"file_path" gorm:"type:varchar(150)"`
	FileType string `json:"file_type" gorm:"type:varchar(10)"`
	Model
}

type FileResponse struct {
	// ID       int    `json:"id"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}
