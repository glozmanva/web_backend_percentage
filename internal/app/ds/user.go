package ds

import "time"

type User struct {
	ID        int       `gorm:"type:integer;primaryKey;autoIncrement"`
	FullName  string    `gorm:"column:full_name;type:varchar(255);not null"`
	BirthDate time.Time `gorm:"column:birth_date;type:date;not null"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
}

func (User) TableName() string {
	return "users"
}
