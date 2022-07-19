package contfileupload

import (
	"context"
	"net/http"

	ifileupload "gitlab.com/369-engineer/369backend/account/interface/fileupload"
	"gitlab.com/369-engineer/369backend/account/models"
	"gitlab.com/369-engineer/369backend/account/pkg/app"
	"gitlab.com/369-engineer/369backend/account/pkg/middleware"
	s3gateway "gitlab.com/369-engineer/369backend/account/pkg/s3"
	util "gitlab.com/369-engineer/369backend/account/pkg/utils"

	"gitlab.com/369-engineer/369backend/account/pkg/logging"
	tool "gitlab.com/369-engineer/369backend/account/pkg/tools"

	"github.com/gin-gonic/gin"
)

// ContFileUpload :
type ContFileUpload struct {
	useSaFileUpload ifileupload.UseCase
}

// NewContFileUpload :
func NewContFileUpload(e *gin.Engine, useSaFileUpload ifileupload.UseCase) {
	cont := &ContFileUpload{
		useSaFileUpload: useSaFileUpload,
	}

	e.Static("/wwwroot", "wwwroot")
	r := e.Group("/v1/account/fileupload")
	// Configure middleware with custom claims
	// r.Use(midd.Versioning)
	r.Use(middleware.Authorize())
	r.POST("", cont.CreateImage)
	r.DELETE("", cont.Delete)

}

// CreateImage :
// @Summary File Upload
// @Security ApiKeyAuth
// @Description Upload file
// @Tags FileUpload
// @Accept  multipart/form-data
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Language header string true "Language Apps"
// @Param upload_file formData file true "account image"
// @Param path formData string true "path images"
// @Success 200 {object} app.Response
// @Router /v1/account/fileupload [post]
func (u *ContFileUpload) CreateImage(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var (
		appE          = app.Gin{C: e}
		imageFormList = []models.FileResponse{}
		logger        = logging.Logger{}
	)

	form, err := e.MultipartForm()
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	images := form.File["upload_file"]

	pt := form.Value["path"]

	logger.Info(pt)
	for _, image := range images {

		S3Output := make(chan *models.FileResponse, 1)
		defer close(S3Output)
		S3 := &s3gateway.DoS3{
			FileUpload: image,
			FileName:   image.Filename,
			PathFile:   pt[0],
		}
		// S3.DeleteFileS3()
		go S3.UploadFileS3(S3Output)
		if err != nil {
			appE.ResponseError(tool.GetStatusCode(err), err)
			return
		}
		var dt = <-S3Output
		imageFormList = append(imageFormList, models.FileResponse{
			FileName: dt.FileName,
			FilePath: dt.FilePath,
			FileType: dt.FileType,
		})

	}
	appE.Response(http.StatusOK, "Ok", imageFormList)

}

// DeleteImages :
// @Summary Delete FileUpload
// @Security ApiKeyAuth
// @Tags FileUpload
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Language header string true "Language Apps"
// @Param req body models.FileResponse true "req param #pakai yang file_path ajah"
// @Success 200 {object} app.Response
// @Router /v1/account/fileupload [delete]
func (u *ContFileUpload) Delete(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.FileResponse{}
	)
	// ID, err := strconv.Atoi(id)
	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}
	S3 := &s3gateway.DoS3{
		FileUpload: nil,
		FileName:   form.FileName,
		PathFile:   form.FilePath,
	}

	S3.DeleteFileS3()

	appE.Response(http.StatusOK, "Ok", nil)
}
