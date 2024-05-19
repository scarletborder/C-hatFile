package models

type User struct {
	ID           uint64 `gorm:"primary_key;auto_increment;"`
	Username     string
	Enc_password string // 一次sha256

	Level uint8 // 用户等级0-匿名,1-管理员，2-admin

	Dirty bool `gorm:"-"` // 忽略这个字段，不写入数据库
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) GetFeature() string {
	return u.Username
}

func (u *User) IsDirty() bool {
	return u.Dirty
}

func (u *User) SetDirty() {
	u.Dirty = true
}

func (u *User) FlushDirty() {
	u.Dirty = false
}
