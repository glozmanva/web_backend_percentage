package ds

type DepositMonthLike struct {
	ID int `gorm:"type:integer;primaryKey;autoIncrement"`

	UserID int `gorm:"column:user_id;type:integer;not null;uniqueIndex:idx_user_deposit_month"`

	DepositMonthID int `gorm:"column:deposit_month_id;type:integer;not null;uniqueIndex:idx_user_deposit_month"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`

	DepositMonth DepositMonth `gorm:"foreignKey:DepositMonthID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (DepositMonthLike) TableName() string {
	return "deposit_month_likes"
}
