package repofileupload

import (
	"context"
	"fmt"

	ifileupload "gitlab.com/369-engineer/369backend/account/interface/fileupload"
	"gitlab.com/369-engineer/369backend/account/models"
	"gitlab.com/369-engineer/369backend/account/pkg/logging"

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
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	// err = db.Conn.Create(userData).Error
	if err != nil {
		return err
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
		//
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
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error

	if err != nil {
		return err
	}
	return nil
}
