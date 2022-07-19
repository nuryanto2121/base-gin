package usesafileupload

import (
	"context"
	ifileupload "gitlab.com/369-engineer/369backend/account/interface/fileupload"
	"gitlab.com/369-engineer/369backend/account/models"
	"time"
)

type useSaFileUpload struct {
	repoSaFileUpload ifileupload.Repository
	contextTimeOut   time.Duration
}

func NewSaFileUpload(a ifileupload.Repository, timeout time.Duration) ifileupload.UseCase {
	return &useSaFileUpload{
		repoSaFileUpload: a,
		contextTimeOut:   timeout,
	}
}

func (u *useSaFileUpload) CreateFileUpload(ctx context.Context, data *models.FileUpload) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		err error
	)

	err = u.repoSaFileUpload.CreateFileUpload(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
func (u *useSaFileUpload) GetBySaFileUpload(ctx context.Context, fileID int) (models.FileUpload, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		err    error
		result models.FileUpload
	)
	result, err = u.GetBySaFileUpload(ctx, fileID)
	if err != nil {
		return result, err
	}
	return result, nil
}
