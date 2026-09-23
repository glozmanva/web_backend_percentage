package ds

import "time"

const (
	DepositMonthStatusDraft     = "draft"
	DepositMonthStatusPublished = "published"
	DepositMonthStatusDeleted   = "deleted"
)

type DepositMonth struct {
	ID int `gorm:"type:integer;primaryKey;autoIncrement"`

	Name             string  `gorm:"type:varchar(50);not null"`
	ShortDescription *string `gorm:"column:short_description;type:varchar(255)"`
	Description      *string `gorm:"type:varchar(1000)"`

	Status string `gorm:"type:varchar(20);not null"`

	ImageURL *string `gorm:"column:image_url;type:varchar(500)"`
	VideoURL *string `gorm:"column:video_url;type:varchar(500)"`

	MonthNumber *int16 `gorm:"column:month_number;type:smallint"`
	DaysCount   *int16 `gorm:"column:days_count;type:smallint"`

	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp;not null;autoCreateTime"`
	FormedAt  *time.Time `gorm:"column:formed_at;type:timestamp"`

	CreatorID int `gorm:"column:creator_id;type:integer;not null"`

	Creator User `gorm:"foreignKey:CreatorID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (DepositMonth) TableName() string {
	return "deposit_months"
}
