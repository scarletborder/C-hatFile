package models

import (
	"fmt"
	"time"
)

type MetaData struct {
	ID         uint `gorm:"primaryKey;auto_increment;"`
	Name       string
	Size       int64
	Tags       []Tag `gorm:"many2many:metadata_tags;"`
	UserID     int64 // 上传用户qq号或其他形式的id
	UploadTime *time.Time
}

/*
@param meta 文件元数据

@return 返回对象的存储名称
*/
func (m MetaData) GenerateObjectName() string {
	return fmt.Sprint(m.ID) + "_" + m.Name
}

func (m MetaData) GetID() uint64 {
	return uint64(m.ID)
}
