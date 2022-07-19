package s3gateway

import (
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3Conn *s3.S3
var S3Session *session.Session

type DoS3 struct {
	FileUpload *multipart.FileHeader `json:"file_upload"`
	FileName   string                `json:"file_name"`
	FileType   string                `json:"file_type"`
	PathFile   string                `json:"path_file"`
}

func Setup() {
	now := time.Now()

	key := setting.S3Setting.SpaceKey
	secret := setting.S3Setting.SpaceSecret

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		// Endpoint:    &setting.S3Setting.EndPoint,
		Region: &setting.S3Setting.Region,
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Printf("dogateway.setup err : %v", err)
		panic(err)
	}
	S3Session = newSession
	S3Conn = s3.New(newSession)

	timeSpent := time.Since(now)
	log.Printf("Config s3 AWS is ready in %v", timeSpent)
}

func (d *DoS3) UploadToDO3() (url string, err error) {
	var tr = true
	src, err := d.FileUpload.Open()
	if err != nil {
		return url, err
	}
	object := s3.PutObjectInput{
		ACL:              aws.String("public-read"),
		Body:             src,
		Bucket:           &setting.S3Setting.SpaceBucket,
		BucketKeyEnabled: &tr,
		Key:              &d.FileName,
		Metadata: map[string]*string{
			"x-amz-meta-my-key": &setting.S3Setting.SpaceKey,
		},
		// ContentType: ,
		Tagging: aws.String("key1=value1&key2=value2"),
	}
	data, err := S3Conn.PutObject(&object)
	if err != nil {
		return url, err
	}
	fmt.Printf("\nreturn DO => %s\n", data.String())
	return url, nil
}
func (d *DoS3) UploadFileS3(output chan *models.FileResponse) (err error) {
	var (
		logger = logging.Logger{}
		tr     = true
		out    = &models.FileResponse{}
		key    = fmt.Sprintf("%s/%s", d.PathFile, d.FileName)
	)
	contentType := genContentType(d)
	src, err := d.FileUpload.Open()
	if err != nil {
		return err
	}
	object := s3.PutObjectInput{
		Body:             src,
		Bucket:           &setting.S3Setting.SpaceBucket,
		BucketKeyEnabled: &tr,
		Key:              &key,
		ContentType:      &contentType,
		Tagging:          aws.String("public=yes"),
		Metadata: map[string]*string{
			"x-amz-meta-my-key": &setting.S3Setting.SpaceKey,
		},
	}

	_, err = S3Conn.PutObject(&object)
	if err != nil {
		logger.Error(err)
		return err
	}

	out.FilePath = setting.S3Setting.EndPoint + key
	out.FileType = contentType
	out.FileName = d.FileName

	output <- out

	return nil
}
func (d *DoS3) GetUrlS3() (url string, err error) {
	var (
		logger = logging.Logger{}
	)
	contentType := genContentType(d)

	req, opt := S3Conn.GetObjectRequest(&s3.GetObjectInput{
		Bucket:              &setting.S3Setting.SpaceBucket,
		Key:                 &d.FileName,
		ResponseContentType: aws.String(contentType),
	})

	url, err = req.Presign(43800 * time.Hour) //5 tahun
	if err != nil {
		logger.Error(err)
		return url, err
	}
	i := strings.Index(url, "?")

	url = url[:i]
	fmt.Println(req)
	fmt.Println(opt)

	return url, nil
}
func (d *DoS3) GetFileS3() (result *s3.GetObjectOutput, err error) {
	var (
		logger = logging.Logger{}
	)
	contentType := genContentType(d)

	req, err := S3Conn.GetObject(&s3.GetObjectInput{
		Bucket:              &setting.S3Setting.SpaceBucket,
		Key:                 &d.FileName,
		ResponseContentType: aws.String(contentType),
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	fmt.Println(req)

	return req, nil
}
func (d *DoS3) DeleteFileS3() (err error) {
	var (
		logger = logging.Logger{}
		key    = strings.Split(d.PathFile, setting.S3Setting.EndPoint) //fmt.Sprintf("%s/%s", d.PathFile, d.FileName)
	)
	_, err = S3Conn.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &setting.S3Setting.SpaceBucket,
		Key:    &key[1],
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
func genContentType(d *DoS3) (result string) {
	extFile := filepath.Ext(d.FileName)
	extImages := setting.AppSetting.ImageAllowExts
	switch extFile {
	case contains(extFile, extImages):
		ext := strings.Replace(extFile, ".", "", 1)
		result = fmt.Sprintf("image/%s", ext)
	case ".pdf":
		result = fmt.Sprintf("application/%s", extFile)
	default:
		result = "application/octet-stream"
	}
	return result
}

func contains(v string, a []string) string {
	for _, i := range a {
		if i == v {
			return i
		}
	}
	return "unknown"
}
