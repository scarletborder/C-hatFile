package chats3

import (
	"chatFileBackend/models"
	"errors"
	"fmt"
	"io"

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
