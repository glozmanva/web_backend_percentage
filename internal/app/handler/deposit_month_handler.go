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
	var depositMonthID int

	idStr := ctx.Query("deposit_month_id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "некорректный id карточки",
				},
			)
			return
		}

		depositMonthID = id
	}

	next := ctx.Query("next") == "true"

	depositMonth, err :=
		h.Repository.GetFeedDepositMonth(
			depositMonthID,
			next,
		)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{
					"error": "карточка не найдена",
				},
			)
			return
		}

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
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

	minValue := ctx.Query("min_days_count")

	if minValue != "" {
		value, err := strconv.Atoi(minValue)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "некорректное минимальное количество дней",
				},
			)
			return
		}

		minDaysCount = value
	}

	maxValue := ctx.Query("max_days_count")

	if maxValue != "" {
		value, err := strconv.Atoi(maxValue)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "некорректное максимальное количество дней",
				},
			)
			return
		}

		maxDaysCount = value
	}

	if minDaysCount < 28 ||
		maxDaysCount > 31 ||
		minDaysCount > maxDaysCount {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "некорректный диапазон количества дней",
			},
		)
		return
	}

	depositMonths, err :=
		h.Repository.GetPublishedDepositMonths(
			minDaysCount,
			maxDaysCount,
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	minDaysPosition :=
		float64(minDaysCount-28) /
			float64(31-28) *
			100

	maxDaysPosition :=
		float64(maxDaysCount-28) /
			float64(31-28) *
			100

	ctx.HTML(
		http.StatusOK,
		"deposit_month_list.html",
		gin.H{
			"deposit_months":    depositMonths,
			"min_days_count":    minDaysCount,
			"max_days_count":    maxDaysCount,
			"min_days_position": minDaysPosition,
			"max_days_position": maxDaysPosition,
		},
	)
}

func (h *Handler) GetDraftDepositMonth(ctx *gin.Context) {
	depositMonth, err :=
		h.Repository.GetDraftDepositMonth(
			currentUserID,
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if depositMonth == nil {
		ctx.HTML(
			http.StatusOK,
			"deposit_month_draft.html",
			gin.H{
				"has_draft": false,
			},
		)
		return
	}

	ctx.HTML(
		http.StatusOK,
		"deposit_month_draft.html",
		gin.H{
			"has_draft":     true,
			"deposit_month": depositMonth,
		},
	)
}

func (h *Handler) CreateDraftDepositMonth(ctx *gin.Context) {
	name := ctx.PostForm("deposit_month_name")

	if name == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "необходимо указать название",
			},
		)
		return
	}

	existingDraft, err :=
		h.Repository.GetDraftDepositMonth(
			currentUserID,
		)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
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
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/deposit_month/draft",
	)
}

func (h *Handler) PublishDepositMonth(ctx *gin.Context) {
	name :=
		ctx.PostForm(
			"deposit_month_name",
		)

	shortDescription :=
		ctx.PostForm(
			"deposit_month_short_description",
		)

	description :=
		ctx.PostForm(
			"deposit_month_description",
		)

	if name == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "необходимо указать название",
			},
		)
		return
	}

	if shortDescription == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "необходимо заполнить краткую информацию",
			},
		)
		return
	}

	monthNumber, err :=
		strconv.Atoi(
			ctx.PostForm(
				"deposit_month_number",
			),
		)

	if err != nil ||
		monthNumber < 1 ||
		monthNumber > 12 {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "номер месяца должен быть от 1 до 12",
			},
		)
		return
	}

	daysCount, err :=
		strconv.Atoi(
			ctx.PostForm(
				"deposit_month_days_count",
			),
		)

	if err != nil ||
		daysCount < 28 ||
		daysCount > 31 {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "количество дней должно быть от 28 до 31",
			},
		)
		return
	}

	err = h.Repository.PublishDraft(
		currentUserID,
		name,
		shortDescription,
		description,
		int16(monthNumber),
		int16(daysCount),
	)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/deposit_month/list",
	)
}

func (h *Handler) DeleteDepositMonth(ctx *gin.Context) {
	idStr :=
		ctx.PostForm(
			"deposit_month_id",
		)

	depositMonthID, err :=
		strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "некорректный id карточки",
			},
		)
		return
	}

	err =
		h.Repository.DeleteDepositMonthSQL(
			depositMonthID,
		)

	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "карточка не найдена",
			},
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/deposit_month/list",
	)
}
