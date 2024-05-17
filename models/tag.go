package models

import "fmt"

type Tag struct {
	ID    uint `gorm:"primaryKey;auto_increment;"`
	Title string

	dirty bool `gorm:"-"` // 忽略这个字段，不写入数据库
}

func (t *Tag) GetID() uint64 {
	return uint64(t.ID)
}

func (t *Tag) GetFeature() string {
	return fmt.Sprint(t.GetID())
}

func (t *Tag) IsDirty() bool {
	return t.dirty
}

func (t *Tag) Dirty() {
	t.dirty = true
}
