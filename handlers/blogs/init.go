package blogs

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type FileInfo struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
}

var errMsg = ""
var files []FileInfo
var fileIndexMap map[string]int

const directory = "blog_files"

func init() {
	daemon := time.NewTicker(10 * time.Minute)
	err := loadFiles(directory)
	if err != nil {
		errMsg = fmt.Sprintf("Fail to read directory %s \n Time: %s",
			err.Error(), time.Now().Format("2006-01-02 15:04:05"))
		logrus.Errorf(errMsg)
	} else {
		errMsg = ""
	}

	go func() {
		for range daemon.C {
			err := loadFiles(directory)
			if err != nil {
				errMsg = fmt.Sprintf("Fail to read directory %s \n Time: %s",
					err.Error(), time.Now().Format("2006-01-02 15:04:05"))
				logrus.Errorf(errMsg)
			} else {
				errMsg = ""
			}
		}
	}()
}
