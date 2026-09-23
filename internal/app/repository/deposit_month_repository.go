package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"deposit_month/internal/app/ds"

	"gorm.io/gorm"
)

type DepositMonthWithLikes struct {
	ds.DepositMonth
	LikesCount int64 `gorm:"column:likes_count"`
}

func (r *Repository) GetFeedDepositMonth(
	depositMonthID int,
	next bool,
) (*DepositMonthWithLikes, error) {

	if depositMonthID == 0 {
		var result DepositMonthWithLikes

		tx := r.db.
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
			Order("deposit_months.id ASC").
			Limit(1).
			Scan(&result)

		if tx.Error != nil {
			return nil, tx.Error
		}

		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}

		return &result, nil
	}

	if !next {
		var result DepositMonthWithLikes

		tx := r.db.
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
				"deposit_months.id = ?",
				depositMonthID,
			).
			Limit(1).
			Scan(&result)

		if tx.Error != nil {
			return nil, tx.Error
		}

		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}

		return &result, nil
	}

	var nextResult DepositMonthWithLikes

	tx := r.db.
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
			"deposit_months.id > ?",
			depositMonthID,
		).
		Order("deposit_months.id ASC").
		Limit(1).
		Scan(&nextResult)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected > 0 {
		return &nextResult, nil
	}

	var firstResult DepositMonthWithLikes

	tx = r.db.
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
		Order("deposit_months.id ASC").
		Limit(1).
		Scan(&firstResult)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &firstResult, nil
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

	err := r.db.Create(&depositMonth).Error

	if err != nil {
		return nil, err
	}

	return &depositMonth, nil
}

func (r *Repository) PublishDraft(
	creatorID int,
	name string,
	shortDescription string,
	description string,
	monthNumber int16,
	daysCount int16,
) error {

	formedAt := time.Now()

	updates := map[string]interface{}{
		"name":              name,
		"short_description": shortDescription,
		"description":       description,
		"month_number":      monthNumber,
		"days_count":        daysCount,
		"status":            ds.DepositMonthStatusPublished,
		"formed_at":         formedAt,
	}

	result := r.db.
		Model(&ds.DepositMonth{}).
		Where(
			"creator_id = ? AND status = ?",
			creatorID,
			ds.DepositMonthStatusDraft,
		).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(
			"черновик пользователя не найден",
		)
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
		`UPDATE deposit_months SET status = $1 WHERE id = $2 AND status = $3`,
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
