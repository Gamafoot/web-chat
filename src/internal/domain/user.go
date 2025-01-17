package domain

type User struct {
	Id       int    `gorm:"primaryKey"`
	Login    string `gorm:"type:varchar(13)"`
	Password string `gorm:"type:text"`
}

func (u User) TableName() string {
	return "users"
}
