package models

type User struct {
	ID           uint64 `gorm:"primary_key;auto_increment;"`
	Username     string
	Enc_password string // 一次sha256
}

func (u User) GetID() uint64 {
	return u.ID
}
