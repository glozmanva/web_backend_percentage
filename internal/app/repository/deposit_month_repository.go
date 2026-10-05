package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"deposit_month/internal/app/ds"

	"github.com/minio/minio-go/v7"
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
			return DepositMonthWithLikes{},
				gorm.ErrRecordNotFound
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
			return DepositMonthWithLikes{},
				gorm.ErrRecordNotFound
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
		return DepositMonthWithLikes{},
			gorm.ErrRecordNotFound
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

// -------------------- API ЛР3 --------------------

func (r *Repository) serializeDepositMonth(
	depositMonth ds.DepositMonth,
	currentUserID int,
) (ds.DepositMonthSerializer, error) {
	var likesCount int64

	err := r.db.
		Model(&ds.DepositMonthLike{}).
		Where(
			"deposit_month_id = ?",
			depositMonth.ID,
		).
		Count(&likesCount).Error

	if err != nil {
		return ds.DepositMonthSerializer{}, err
	}

	var currentUserLikeCount int64

	err = r.db.
		Model(&ds.DepositMonthLike{}).
		Where(
			"deposit_month_id = ? AND user_id = ?",
			depositMonth.ID,
			currentUserID,
		).
		Count(&currentUserLikeCount).Error

	if err != nil {
		return ds.DepositMonthSerializer{}, err
	}

	isLiked := 0
	if currentUserLikeCount > 0 {
		isLiked = 1
	}

	isCreator := 0
	if depositMonth.CreatorID == currentUserID {
		isCreator = 1
	}

	return ds.DepositMonthSerializer{
		DepositMonth: depositMonth,
		LikesCount:   likesCount,
		IsLiked:      isLiked,
		IsCreator:    isCreator,
	}, nil
}

func (r *Repository) GetPublishedDepositMonthsAPI(
	minDaysCount int,
	maxDaysCount int,
	currentUserID int,
) ([]ds.DepositMonthSerializer, error) {
	var depositMonths []ds.DepositMonth

	err := r.db.
		Where(
			"status = ?",
			ds.DepositMonthStatusPublished,
		).
		Where(
			"days_count BETWEEN ? AND ?",
			minDaysCount,
			maxDaysCount,
		).
		Order("month_number ASC").
		Find(&depositMonths).Error

	if err != nil {
		return nil, err
	}

	result := make(
		[]ds.DepositMonthSerializer,
		0,
		len(depositMonths),
	)

	for _, depositMonth := range depositMonths {
		serialized, err := r.serializeDepositMonth(
			depositMonth,
			currentUserID,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, serialized)
	}

	return result, nil
}

func (r *Repository) GetFeedDepositMonthAPI(
	depositMonthID int,
	next bool,
	currentUserID int,
) (ds.DepositMonthSerializer, error) {
	var depositMonth ds.DepositMonth

	if next && depositMonthID > 0 {
		err := r.db.
			Where(
				"status = ? AND id > ?",
				ds.DepositMonthStatusPublished,
				depositMonthID,
			).
			Order("id ASC").
			First(&depositMonth).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = r.db.
				Where(
					"status = ?",
					ds.DepositMonthStatusPublished,
				).
				Order("id ASC").
				First(&depositMonth).Error
		}

		if err != nil {
			return ds.DepositMonthSerializer{}, err
		}

		return r.serializeDepositMonth(
			depositMonth,
			currentUserID,
		)
	}

	if depositMonthID > 0 {
		err := r.db.
			Where(
				"id = ? AND status = ?",
				depositMonthID,
				ds.DepositMonthStatusPublished,
			).
			First(&depositMonth).Error

		if err != nil {
			return ds.DepositMonthSerializer{}, err
		}

		return r.serializeDepositMonth(
			depositMonth,
			currentUserID,
		)
	}

	err := r.db.
		Where(
			"status = ?",
			ds.DepositMonthStatusPublished,
		).
		Order("id ASC").
		First(&depositMonth).Error

	if err != nil {
		return ds.DepositMonthSerializer{}, err
	}

	return r.serializeDepositMonth(
		depositMonth,
		currentUserID,
	)
}

func (r *Repository) DeleteDepositMonthAPI(
	depositMonthID int,
	creatorID int,
) error {
	result := r.db.
		Model(&ds.DepositMonth{}).
		Where(
			"id = ? AND creator_id = ? AND status IN ?",
			depositMonthID,
			creatorID,
			[]string{
				ds.DepositMonthStatusDraft,
				ds.DepositMonthStatusPublished,
			},
		).
		Update(
			"status",
			ds.DepositMonthStatusDeleted,
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) SetDepositMonthLike(
	depositMonthID int,
	userID int,
	value int,
) error {
	var count int64

	err := r.db.
		Model(&ds.DepositMonth{}).
		Where(
			"id = ? AND status = ?",
			depositMonthID,
			ds.DepositMonthStatusPublished,
		).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	if value == 1 {
		like := ds.DepositMonthLike{
			UserID:         userID,
			DepositMonthID: depositMonthID,
		}

		return r.db.
			Where(
				"user_id = ? AND deposit_month_id = ?",
				userID,
				depositMonthID,
			).
			FirstOrCreate(&like).Error
	}

	return r.db.
		Where(
			"user_id = ? AND deposit_month_id = ?",
			userID,
			depositMonthID,
		).
		Delete(&ds.DepositMonthLike{}).Error
}

func (r *Repository) uploadDepositMonthFile(
	depositMonthID int,
	header *multipart.FileHeader,
	kind string,
) error {
	file, err := header.Open()
	if err != nil {
		return fmt.Errorf(
			"ошибка открытия файла: %w",
			err,
		)
	}
	defer file.Close()

	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf(
			"ошибка чтения файла: %w",
			err,
		)
	}

	if n == 0 {
		return fmt.Errorf("пустой файл")
	}

	contentType := http.DetectContentType(
		buffer[:n],
	)

	if kind == "image" &&
		!strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf(
			"файл должен быть изображением",
		)
	}

	if kind == "video" &&
		!strings.HasPrefix(contentType, "video/") {
		return fmt.Errorf(
			"файл должен быть видео",
		)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf(
			"ошибка перемещения по файлу: %w",
			err,
		)
	}

	extensions := map[string]string{
		"image/jpeg":      ".jpg",
		"image/png":       ".png",
		"image/gif":       ".gif",
		"image/webp":      ".webp",
		"video/mp4":       ".mp4",
		"video/webm":      ".webm",
		"video/quicktime": ".mov",
	}

	extension, ok := extensions[contentType]
	if !ok {
		return fmt.Errorf(
			"неподдерживаемый тип файла: %s",
			contentType,
		)
	}

	filename := fmt.Sprintf(
		"deposit_month_%d_%s%s",
		depositMonthID,
		kind,
		extension,
	)

	ctx := context.Background()

	_, err = r.minio.PutObject(
		ctx,
		r.minioBucketName,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"ошибка загрузки файла в MinIO: %w",
			err,
		)
	}

	fileURL := fmt.Sprintf(
		"http://%s/%s/%s",
		r.minioEndpoint,
		r.minioBucketName,
		filename,
	)

	column := "image_url"
	if kind == "video" {
		column = "video_url"
	}

	err = r.db.
		Model(&ds.DepositMonth{}).
		Where(
			"id = ?",
			depositMonthID,
		).
		UpdateColumn(
			column,
			fileURL,
		).Error

	if err != nil {
		_ = r.minio.RemoveObject(
			ctx,
			r.minioBucketName,
			filename,
			minio.RemoveObjectOptions{},
		)

		return fmt.Errorf(
			"ошибка сохранения пути файла: %w",
			err,
		)
	}

	return nil
}

func (r *Repository) AddOrReplaceDepositMonthImage(
	depositMonthID int,
	header *multipart.FileHeader,
) error {
	return r.uploadDepositMonthFile(
		depositMonthID,
		header,
		"image",
	)
}

func (r *Repository) AddOrReplaceDepositMonthVideo(
	depositMonthID int,
	header *multipart.FileHeader,
) error {
	return r.uploadDepositMonthFile(
		depositMonthID,
		header,
		"video",
	)
}
