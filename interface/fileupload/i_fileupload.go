package ifileupload

import (
	"app/models"
	"context"
)

type Repository interface {
	CreateFileUpload(ctx context.Context, data *models.FileUpload) error
	GetBySaFileUpload(ctx context.Context, fileID int) (models.FileUpload, error)
	DeleteSaFileUpload(ctx context.Context, fileID int) error
}

type UseCase interface {
	CreateFileUpload(ctx context.Context, data *models.FileUpload) error
	GetBySaFileUpload(ctx context.Context, fileID int) (models.FileUpload, error)
}
