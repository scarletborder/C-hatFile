package models

type User struct {
	ID           uint64 `gorm:"primary_key;auto_increment;"`
	Username     string
	Enc_password string // 一次sha256

	dirty bool `gorm:"-"` // 忽略这个字段，不写入数据库
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) GetFeature() string {
	return u.Username
}

func (u *User) IsDirty() bool {
	return u.dirty
}

func (u *User) Dirty() {
	u.dirty = true
}
