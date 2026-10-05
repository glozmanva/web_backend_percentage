package ds

type DepositMonthLike struct {
	ID             int          `gorm:"type:integer;primaryKey;autoIncrement" json:"id"`
	UserID         int          `gorm:"column:user_id;type:integer;not null;uniqueIndex:idx_user_deposit_month" json:"user_id"`
	DepositMonthID int          `gorm:"column:deposit_month_id;type:integer;not null;uniqueIndex:idx_user_deposit_month" json:"deposit_month_id"`
	User           User         `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;" json:"-"`
	DepositMonth   DepositMonth `gorm:"foreignKey:DepositMonthID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;" json:"-"`
}

func (DepositMonthLike) TableName() string {
	return "deposit_month_likes"
}
