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

	dirty bool `gorm:"-"` // 忽略这个字段，不写入数据库
}

/*
@param meta 文件元数据

@return 返回对象的存储名称
*/
func (m *MetaData) GenerateObjectName() string {
	return fmt.Sprint(m.ID) + "_" + m.Name
}

func (m *MetaData) GetID() uint64 {
	return uint64(m.ID)
}

func (m *MetaData) GetFeature() string {
	return fmt.Sprint(m.GetID())
}

func (m *MetaData) IsDirty() bool {
	return m.dirty
}

func (m *MetaData) Dirty() {
	m.dirty = true
}
