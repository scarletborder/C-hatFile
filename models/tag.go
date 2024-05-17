package models

type Tag struct {
	ID    uint `gorm:"primaryKey;auto_increment;"`
	Title string
}

func (t Tag) GetID() uint64 {
	return uint64(t.ID)
}
