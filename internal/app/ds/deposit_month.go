package ds

import "time"

const (
	DepositMonthStatusDraft     = "draft"
	DepositMonthStatusPublished = "published"
	DepositMonthStatusDeleted   = "deleted"
)

type DepositMonth struct {
	ID          int        `gorm:"type:integer;primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"type:varchar(50);not null" json:"name"`
	Description *string    `gorm:"type:varchar(1000)" json:"description"`
	Status      string     `gorm:"type:varchar(20);not null" json:"status"`
	ImageURL    *string    `gorm:"column:image_url;type:varchar(500)" json:"image_url"`
	VideoURL    *string    `gorm:"column:video_url;type:varchar(500)" json:"video_url"`
	MonthNumber *int16     `gorm:"column:month_number;type:smallint" json:"month_number"`
	DaysCount   *int16     `gorm:"column:days_count;type:smallint" json:"days_count"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;not null;autoCreateTime" json:"created_at"`
	FormedAt    *time.Time `gorm:"column:formed_at;type:timestamp" json:"formed_at"`
	CreatorID   int        `gorm:"column:creator_id;type:integer;not null" json:"-"`
	Creator     User       `gorm:"foreignKey:CreatorID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;" json:"-"`
}

type DepositMonthSerializer struct {
	DepositMonth
	LikesCount int64 `json:"likes_count"`
	IsLiked    int   `json:"is_liked"`
	IsCreator  int   `json:"is_creator"`
}

func (DepositMonth) TableName() string {
	return "deposit_months"
}
