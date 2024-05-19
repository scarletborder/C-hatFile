package chats3

import (
	"chatFileBackend/models"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type S3Point struct {
	EndPoint     string        `json:"your-s3-endpoint-com"`
	CustomDomain string        `json:"custom-domain" default:""`
	AccessKey    string        `json:"your-access-key"`
	SecretKey    string        `json:"your-secret-key"`
	ChunkSize    uint64        `json:"chunksize" default:"0"`
	UseSSL       bool          `json:"usessl" default:"false"`
	Point        *minio.Client `default:"null"`
}

func (s3p S3Point) ExistBucket(bucket_name string) bool {
	found, err := s3p.Point.BucketExists(context.Background(), bucket_name)
	if err != nil {
		logrus.Errorln(err)
	}
	return found
}

func (s3p S3Point) MakeBucket(bucket_name string) {
	err := s3p.Point.MakeBucket(context.Background(), bucket_name, minio.MakeBucketOptions{}) // Region: "us-east-1"
	if err != nil {
		logrus.Errorln(err)
	}
}

func (s3p S3Point) upload(bucket_name string, meta *models.MetaData, file io.Reader) error {
	_, err := s3p.Point.PutObject(context.Background(), bucket_name,
		meta.GenerateObjectName(),
		file, meta.Size,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return err
}

// 返回下载链接或者错误
func (s3p S3Point) downloadURL(bucket_name string, meta *models.MetaData) (string, error) {
	expires := time.Duration(1) * time.Hour // 链接有效时间，例如1小时

	// 生成预签名的URL
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition",
		fmt.Sprintf("attachment; filename=\"%v\"", meta.Name))
	presignedURL, err := s3p.Point.PresignedGetObject(context.Background(), bucket_name,
		meta.GenerateObjectName(), expires,
		reqParams)
	// 配置自定义域名: 首先，你需要将你的自定义域名指向MinIO服务器。
	// 这通常涉及到在你的DNS提供商处添加一个CNAME记录，将你的自定义域名指向MinIO服务器的地址。
	if err != nil {
		return "", err
	}
	ret := presignedURL.String()
	if s3p.CustomDomain != "" {
		ret = strings.Replace(ret, s3p.EndPoint, s3p.CustomDomain, -1)
	}
	return ret, err
}

// 返回下载流或错误
func (s3p S3Point) DownloadReader(bucket_name string, meta models.MetaData) (io.Reader, error) {
	reader, err := s3p.Point.GetObject(context.Background(), bucket_name,
		meta.GenerateObjectName(), minio.GetObjectOptions{})
	return reader, err
}

func (s3p S3Point) DownloadReaderByObjectName(bucket_name, obj_name string) (io.Reader, error) {
	reader, err := s3p.Point.GetObject(context.Background(), bucket_name,
		obj_name, minio.GetObjectOptions{})
	return reader, err
}
