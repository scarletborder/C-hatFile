package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type MetaData struct {
	ID   uint `gorm:"primaryKey;auto_increment;"`
	Name string
	Size int64
	Tags []Tag `gorm:"many2many:metadata_tags;"`
	// Username   string // 上传用户qq号或其他形式的id
	UserID     uint64
	UploadTime *time.Time

	Dirty bool `gorm:"-"` // 忽略这个字段，不写入数据库
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
	return m.Dirty
}

func (m *MetaData) SetDirty() {
	m.Dirty = true
}

func (m *MetaData) FlushDirty() {
	m.Dirty = false
}

// json
type MetaDataNameOnly struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Tags       []string   `json:"tags"`
	UploadTime *time.Time `json:"upload_time"`
}

func (m MetaData) MarshalJSON() ([]byte, error) {
	var tags []string

	for _, t := range m.Tags {
		tags = append(tags, t.Title)
	}

	return json.Marshal(MetaDataNameOnly{
		Name:       m.Name,
		Tags:       tags,
		UploadTime: m.UploadTime,
		ID:         m.ID,
	})
}
