package ds

import "time"

type User struct {
	ID        int       `gorm:"type:integer;primaryKey;autoIncrement" json:"id"`
	FullName  string    `gorm:"column:full_name;type:varchar(255);not null" json:"full_name"`
	BirthDate time.Time `gorm:"column:birth_date;type:date;not null" json:"birth_date"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}
