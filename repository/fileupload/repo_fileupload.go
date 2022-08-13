package repofileupload

import (
	"context"
	"fmt"

	ifileupload "app/interface/fileupload"
	"app/models"
	"app/pkg/logging"

	"gorm.io/gorm"
)

type repoAuth struct {
	Conn *gorm.DB
}

func NewRepoFileUpload(Conn *gorm.DB) ifileupload.Repository {
	return &repoAuth{Conn}
}

func (m *repoAuth) CreateFileUpload(ctx context.Context, data *models.FileUpload) (err error) {
	var logger = logging.Logger{}
	query := m.Conn.Create(&data)

	err = query.Error
	// err = db.Conn.Create(userData).Error
	if err != nil {
		logger.Error("repo fileupload CreateFileUpload ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (m *repoAuth) GetBySaFileUpload(ctx context.Context, fileID int) (models.FileUpload, error) {
	var (
		dataFileUpload = models.FileUpload{}
		logger         = logging.Logger{}
		err            error
	)
	query := m.Conn.Where("file_id = ?", fileID).First(&dataFileUpload)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error

	if err != nil {
		logger.Error("repo fileupload GetBySaFileUpload ", err)
		if err == gorm.ErrRecordNotFound {
			return dataFileUpload, models.ErrNotFound
		}
		return dataFileUpload, err
	}

	return dataFileUpload, err
}
func (m *repoAuth) DeleteSaFileUpload(ctx context.Context, fileID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	userData := &models.FileUpload{}
	userData.ID = fileID

	query := m.Conn.Delete(&userData)

	err = query.Error

	if err != nil {
		logger.Error("repo fileupload DeleteSaFileUpload ", err)
		return models.ErrInternalServerError
	}
	return nil
}
