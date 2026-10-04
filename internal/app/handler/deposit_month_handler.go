package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const currentUserID int = 1

func (h *Handler) GetDepositMonth(ctx *gin.Context) {
	depositMonthID := 0

	if idString := ctx.Query("deposit_month_id"); idString != "" {
		id, err := strconv.Atoi(idString)
		if err != nil || id <= 0 {
			ctx.String(
				http.StatusBadRequest,
				"Некорректный id расчётного месяца",
			)
			return
		}

		depositMonthID = id
	}

	next := ctx.Query("next") == "true"

	depositMonth, err := h.Repository.GetFeedDepositMonth(
		depositMonthID,
		next,
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.String(
			http.StatusNotFound,
			"Расчётный месяц не найден",
		)
		return
	}

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.HTML(
		http.StatusOK,
		"deposit_month_feed.html",
		gin.H{
			"deposit_month": depositMonth.DepositMonth,
			"likes_count":   depositMonth.LikesCount,
		},
	)
}

func (h *Handler) GetDepositMonths(ctx *gin.Context) {
	minDaysCount := 28
	maxDaysCount := 31

	if value := ctx.Query("min_days_count"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed >= 28 && parsed <= 31 {
			minDaysCount = parsed
		}
	}

	if value := ctx.Query("max_days_count"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed >= 28 && parsed <= 31 {
			maxDaysCount = parsed
		}
	}

	if minDaysCount > maxDaysCount {
		minDaysCount, maxDaysCount =
			maxDaysCount, minDaysCount
	}

	depositMonths, err :=
		h.Repository.GetPublishedDepositMonths(
			minDaysCount,
			maxDaysCount,
		)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	minPosition :=
		float64(minDaysCount-28) / 3.0 * 100.0

	maxPosition :=
		float64(maxDaysCount-28) / 3.0 * 100.0

	ctx.HTML(
		http.StatusOK,
		"deposit_month_list.html",
		gin.H{
			"deposit_months":    depositMonths,
			"min_days_count":    minDaysCount,
			"max_days_count":    maxDaysCount,
			"min_days_position": minPosition,
			"max_days_position": maxPosition,
		},
	)
}

func (h *Handler) GetDraftDepositMonth(
	ctx *gin.Context,
) {
	depositMonth, err :=
		h.Repository.GetDraftDepositMonth(currentUserID)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.HTML(
		http.StatusOK,
		"deposit_month_draft.html",
		gin.H{
			"has_draft":     depositMonth != nil,
			"deposit_month": depositMonth,
		},
	)
}

func (h *Handler) CreateDraftDepositMonth(
	ctx *gin.Context,
) {
	name := ctx.PostForm("deposit_month_name")

	if name == "" {
		ctx.String(
			http.StatusBadRequest,
			"Название обязательно",
		)
		return
	}

	existingDraft, err :=
		h.Repository.GetDraftDepositMonth(currentUserID)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	if existingDraft != nil {
		ctx.Redirect(
			http.StatusFound,
			"/deposit_month/draft",
		)
		return
	}

	_, err = h.Repository.CreateDraft(
		name,
		currentUserID,
	)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/deposit_month/draft",
	)
}

func (h *Handler) PublishDepositMonth(
	ctx *gin.Context,
) {
	name := ctx.PostForm("deposit_month_name")
	description := ctx.PostForm("deposit_month_description")

	monthNumber, err := strconv.Atoi(
		ctx.PostForm("deposit_month_month_number"),
	)

	if err != nil ||
		monthNumber < 1 ||
		monthNumber > 12 {
		ctx.String(
			http.StatusBadRequest,
			"Некорректный номер месяца",
		)
		return
	}

	daysCount, err := strconv.Atoi(
		ctx.PostForm("deposit_month_days_count"),
	)

	if err != nil ||
		daysCount < 28 ||
		daysCount > 31 {
		ctx.String(
			http.StatusBadRequest,
			"Некорректное количество дней",
		)
		return
	}

	if name == "" {
		ctx.String(
			http.StatusBadRequest,
			"Название обязательно",
		)
		return
	}

	if description == "" {
		ctx.String(
			http.StatusBadRequest,
			"Описание обязательно",
		)
		return
	}

	err = h.Repository.PublishDraft(
		currentUserID,
		name,
		description,
		int16(monthNumber),
		int16(daysCount),
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.String(
			http.StatusNotFound,
			"Черновик не найден",
		)
		return
	}

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/deposit_month/list",
	)
}

func (h *Handler) DeleteDepositMonth(
	ctx *gin.Context,
) {
	id, err := strconv.Atoi(
		ctx.PostForm("deposit_month_id"),
	)

	if err != nil || id <= 0 {
		ctx.String(
			http.StatusBadRequest,
			"Некорректный id",
		)
		return
	}

	err = h.Repository.DeleteDepositMonthSQL(id)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/deposit_month/list",
	)
}
