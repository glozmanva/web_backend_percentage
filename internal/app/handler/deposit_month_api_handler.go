package handler

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"deposit_month/internal/app/currentuser"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func validateDepositMonthFile(
	header *multipart.FileHeader,
	kind string,
) error {
	file, err := header.Open()
	if err != nil {
		return fmt.Errorf(
			"не удалось открыть файл",
		)
	}
	defer file.Close()

	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf(
			"не удалось прочитать файл",
		)
	}

	if n == 0 {
		return fmt.Errorf("файл пустой")
	}

	contentType := http.DetectContentType(
		buffer[:n],
	)

	if kind == "image" {
		allowed := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}

		if !allowed[contentType] {
			return fmt.Errorf(
				"файл image должен быть изображением",
			)
		}
	}

	if kind == "video" {
		allowed := map[string]bool{
			"video/mp4":       true,
			"video/webm":      true,
			"video/quicktime": true,
		}

		if !allowed[contentType] {
			return fmt.Errorf(
				"файл video должен быть видео",
			)
		}
	}

	return nil
}

func (h *Handler) GetDepositMonthsAPI(
	ctx *gin.Context,
) {
	minDaysCount := 28
	maxDaysCount := 31

	if value := ctx.Query("min_days_count"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 28 || parsed > 31 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "некорректный min_days_count",
				},
			)
			return
		}

		minDaysCount = parsed
	}

	if value := ctx.Query("max_days_count"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 28 || parsed > 31 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "некорректный max_days_count",
				},
			)
			return
		}

		maxDaysCount = parsed
	}

	if minDaysCount > maxDaysCount {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "min_days_count больше max_days_count",
			},
		)
		return
	}

	depositMonths, err :=
		h.Repository.GetPublishedDepositMonthsAPI(
			minDaysCount,
			maxDaysCount,
			currentuser.GetCurrentUserID(),
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		depositMonths,
	)
}

func (h *Handler) GetDepositMonthFeedAPI(
	ctx *gin.Context,
) {
	depositMonthID := 0

	if value := ctx.Query("deposit_month_id"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "некорректный deposit_month_id",
				},
			)
			return
		}

		depositMonthID = parsed
	}

	next := false

	if value := ctx.Query("next"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "next должен быть true или false",
				},
			)
			return
		}

		next = parsed
	}

	depositMonth, err :=
		h.Repository.GetFeedDepositMonthAPI(
			depositMonthID,
			next,
			currentuser.GetCurrentUserID(),
		)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "опубликованная запись не найдена",
			},
		)
		return
	}

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		depositMonth,
	)
}

func (h *Handler) GetDraftDepositMonthAPI(
	ctx *gin.Context,
) {
	depositMonth, err :=
		h.Repository.GetDraftDepositMonth(
			currentuser.GetCurrentUserID(),
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	if depositMonth == nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "черновик не найден",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		depositMonth,
	)
}

func (h *Handler) CreateDepositMonthAPI(
	ctx *gin.Context,
) {
	currentUserID :=
		currentuser.GetCurrentUserID()

	existingDraft, err :=
		h.Repository.GetDraftDepositMonth(
			currentUserID,
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	if existingDraft != nil {
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"error": "у пользователя уже есть черновик",
			},
		)
		return
	}

	err = ctx.Request.ParseMultipartForm(
		64 << 20,
	)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "некорректная multipart-форма",
			},
		)
		return
	}

	name := ctx.Request.FormValue("name")

	if name == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "поле name обязательно",
			},
		)
		return
	}

	imageHeader, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "файл image обязателен",
			},
		)
		return
	}

	videoHeader, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "файл video обязателен",
			},
		)
		return
	}

	if err := validateDepositMonthFile(
		imageHeader,
		"image",
	); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := validateDepositMonthFile(
		videoHeader,
		"video",
	); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	depositMonth, err :=
		h.Repository.CreateDraft(
			name,
			currentUserID,
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = h.Repository.AddOrReplaceDepositMonthImage(
		depositMonth.ID,
		imageHeader,
	)

	if err != nil {
		_ = h.Repository.DeleteDepositMonthAPI(
			depositMonth.ID,
			currentUserID,
		)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = h.Repository.AddOrReplaceDepositMonthVideo(
		depositMonth.ID,
		videoHeader,
	)

	if err != nil {
		_ = h.Repository.DeleteDepositMonthAPI(
			depositMonth.ID,
			currentUserID,
		)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	depositMonth, err =
		h.Repository.GetDraftDepositMonth(
			currentUserID,
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		depositMonth,
	)
}

func (h *Handler) PublishDepositMonthAPI(
	ctx *gin.Context,
) {
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		MonthNumber int16  `json:"month_number"`
		DaysCount   int16  `json:"days_count"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "некорректное тело запроса",
			},
		)
		return
	}

	if request.Name == "" ||
		request.Description == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "name и description обязательны",
			},
		)
		return
	}

	if request.MonthNumber < 1 ||
		request.MonthNumber > 12 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "month_number должен быть от 1 до 12",
			},
		)
		return
	}

	if request.DaysCount < 28 ||
		request.DaysCount > 31 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "days_count должен быть от 28 до 31",
			},
		)
		return
	}

	err := h.Repository.PublishDraft(
		currentuser.GetCurrentUserID(),
		request.Name,
		request.Description,
		request.MonthNumber,
		request.DaysCount,
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "черновик не найден",
			},
		)
		return
	}

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "услуга опубликована",
		},
	)
}

func (h *Handler) DeleteDepositMonthAPI(
	ctx *gin.Context,
) {
	id, err := strconv.Atoi(
		ctx.Param("id"),
	)

	if err != nil || id <= 0 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "некорректный id",
			},
		)
		return
	}

	err = h.Repository.DeleteDepositMonthAPI(
		id,
		currentuser.GetCurrentUserID(),
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "запись не найдена или принадлежит другому пользователю",
			},
		)
		return
	}

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "услуга удалена",
		},
	)
}

func (h *Handler) SetDepositMonthLikeAPI(
	ctx *gin.Context,
) {
	id, err := strconv.Atoi(
		ctx.Param("id"),
	)

	if err != nil || id <= 0 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "некорректный id",
			},
		)
		return
	}

	var request struct {
		Value *int `json:"value" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "поле value обязательно",
			},
		)
		return
	}

	if *request.Value != 0 &&
		*request.Value != 1 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "value может быть только 0 или 1",
			},
		)
		return
	}

	err = h.Repository.SetDepositMonthLike(
		id,
		currentuser.GetCurrentUserID(),
		*request.Value,
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "опубликованная услуга не найдена",
			},
		)
		return
	}

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"value": *request.Value,
		},
	)
}
