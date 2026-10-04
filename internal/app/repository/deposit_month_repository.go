package repository

import (
	"database/sql"
	"errors"
	"time"

	"deposit_month/internal/app/ds"

	"gorm.io/gorm"
)

type DepositMonthWithLikes struct {
	ds.DepositMonth
	LikesCount int64 `gorm:"column:likes_count"`
}

func (r *Repository) feedQuery() *gorm.DB {
	return r.db.
		Model(&ds.DepositMonth{}).
		Select(`
			deposit_months.*,
			(
				SELECT COUNT(*)
				FROM deposit_month_likes
				WHERE deposit_month_likes.deposit_month_id = deposit_months.id
			) AS likes_count
		`).
		Where(
			"deposit_months.status = ?",
			ds.DepositMonthStatusPublished,
		)
}

func (r *Repository) GetFeedDepositMonth(
	depositMonthID int,
	next bool,
) (DepositMonthWithLikes, error) {
	var depositMonth DepositMonthWithLikes

	if next && depositMonthID > 0 {
		result := r.feedQuery().
			Where("deposit_months.id > ?", depositMonthID).
			Order("deposit_months.id ASC").
			Limit(1).
			Scan(&depositMonth)

		if result.Error != nil {
			return DepositMonthWithLikes{}, result.Error
		}

		if result.RowsAffected == 1 {
			return depositMonth, nil
		}

		depositMonth = DepositMonthWithLikes{}

		result = r.feedQuery().
			Order("deposit_months.id ASC").
			Limit(1).
			Scan(&depositMonth)

		if result.Error != nil {
			return DepositMonthWithLikes{}, result.Error
		}

		if result.RowsAffected == 0 {
			return DepositMonthWithLikes{}, gorm.ErrRecordNotFound
		}

		return depositMonth, nil
	}

	if depositMonthID > 0 {
		result := r.feedQuery().
			Where("deposit_months.id = ?", depositMonthID).
			Limit(1).
			Scan(&depositMonth)

		if result.Error != nil {
			return DepositMonthWithLikes{}, result.Error
		}

		if result.RowsAffected == 0 {
			return DepositMonthWithLikes{}, gorm.ErrRecordNotFound
		}

		return depositMonth, nil
	}

	result := r.feedQuery().
		Order("deposit_months.id ASC").
		Limit(1).
		Scan(&depositMonth)

	if result.Error != nil {
		return DepositMonthWithLikes{}, result.Error
	}

	if result.RowsAffected == 0 {
		return DepositMonthWithLikes{}, gorm.ErrRecordNotFound
	}

	return depositMonth, nil
}

func (r *Repository) GetPublishedDepositMonths(
	minDaysCount int,
	maxDaysCount int,
) ([]DepositMonthWithLikes, error) {
	var depositMonths []DepositMonthWithLikes

	err := r.db.
		Model(&ds.DepositMonth{}).
		Select(`
			deposit_months.*,
			(
				SELECT COUNT(*)
				FROM deposit_month_likes
				WHERE deposit_month_likes.deposit_month_id = deposit_months.id
			) AS likes_count
		`).
		Where(
			"deposit_months.status = ?",
			ds.DepositMonthStatusPublished,
		).
		Where(
			"deposit_months.days_count BETWEEN ? AND ?",
			minDaysCount,
			maxDaysCount,
		).
		Order("deposit_months.month_number ASC").
		Scan(&depositMonths).Error

	if err != nil {
		return nil, err
	}

	return depositMonths, nil
}

func (r *Repository) GetDraftDepositMonth(
	creatorID int,
) (*ds.DepositMonth, error) {
	var depositMonth ds.DepositMonth

	err := r.db.
		Where(
			"creator_id = ? AND status = ?",
			creatorID,
			ds.DepositMonthStatusDraft,
		).
		First(&depositMonth).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &depositMonth, nil
}

func (r *Repository) CreateDraft(
	name string,
	creatorID int,
) (*ds.DepositMonth, error) {
	depositMonth := ds.DepositMonth{
		Name:      name,
		Status:    ds.DepositMonthStatusDraft,
		CreatorID: creatorID,
	}

	if err := r.db.Create(&depositMonth).Error; err != nil {
		return nil, err
	}

	return &depositMonth, nil
}

func (r *Repository) PublishDraft(
	creatorID int,
	name string,
	description string,
	monthNumber int16,
	daysCount int16,
) error {
	now := time.Now()

	result := r.db.
		Model(&ds.DepositMonth{}).
		Where(
			"creator_id = ? AND status = ?",
			creatorID,
			ds.DepositMonthStatusDraft,
		).
		Updates(map[string]interface{}{
			"name":         name,
			"description":  description,
			"month_number": monthNumber,
			"days_count":   daysCount,
			"status":       ds.DepositMonthStatusPublished,
			"formed_at":    &now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) DeleteDepositMonthSQL(
	depositMonthID int,
) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	result, err := sqlDB.Exec(
		`
		UPDATE deposit_months
		SET status = $1
		WHERE id = $2
		  AND status = $3
		`,
		ds.DepositMonthStatusDeleted,
		depositMonthID,
		ds.DepositMonthStatusPublished,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
