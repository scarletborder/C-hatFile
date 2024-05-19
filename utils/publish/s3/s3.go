package chats3

import (
	"chatFileBackend/models"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

var (
	s3points    []S3Point
	bucket_name string
)

func Upload_file(reader io.Reader, meta *models.MetaData) (string, error) {
	for _, s3p := range s3points {
		err := s3p.upload(bucket_name, meta, reader)
		if err == nil {
			return fmt.Sprintln("Successfully upload to " + s3p.EndPoint), nil
		}
		logrus.Errorf("failed upload %v :%v", s3p.EndPoint, err.Error())
	}
	return "", errors.New("all s3points failed to upload file")
}

func Get_download_url(meta *models.MetaData) (string, error) {
	for _, s3p := range s3points {
		url, err := s3p.downloadURL(bucket_name, meta)
		if err == nil {
			return url, nil
		}
		logrus.Errorf("failed download from %v :%v", s3p.EndPoint, err.Error())
	}
	return "All s3points failed to download file", errors.New("all s3points failed to download file")
}

func GetDownlodReader(meta *models.MetaData) (io.Reader, int64, error) {
	return GetDownlodReaderByObjectName(meta.GenerateObjectName())
}

func GetDownlodReaderByObjectName(obj_name string) (io.Reader, int64, error) {
	for _, s3p := range s3points {
		r, err := s3p.DownloadReaderByObjectName(bucket_name, obj_name)
		if err == nil {
			obj_info, _ := s3p.Point.StatObject(context.Background(),
				bucket_name, obj_name,
				minio.StatObjectOptions{})

			return r, obj_info.Size, nil
		}
	}
	return nil, 0, errors.New("no such file " + obj_name)
}

// func Get_download_url_by_fnv(fnv, file_name string) (string, error) {
// 	for _, s3p := range s3points {
// 		url, err := s3p.downloadURL_by_fnv("bot", fnv, file_name)
// 		if err == nil {
// 			return url, nil
// 		}
// 		logrus.Errorf("failed download from %v :%v", s3p.EndPoint, err.Error())
// 	}
// 	return "All s3points failed to download file", errors.New("all s3points failed to download file")
// }
