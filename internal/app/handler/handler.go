package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"web_backend_percentage/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

type MonthListItem struct {
	ID         int
	Name       string
	DaysCount  int
	ImageURL   string
	VideoURL   string
	LikesCount int
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetDepositMonths(ctx *gin.Context) {
	daysString := strings.TrimSpace(ctx.Query("days"))

	var daysFilter *int

	if daysString != "" {
		days, err := strconv.Atoi(daysString)

		if err != nil {
			ctx.String(
				http.StatusBadRequest,
				"Количество дней должно быть числом",
			)
			return
		}

		daysFilter = &days
	}

	months := h.Repository.GetPublishedDepositMonths(daysFilter)

	items := make(
		[]MonthListItem,
		0,
		len(months),
	)

	for _, month := range months {
		items = append(
			items,
			MonthListItem{
				ID:         month.ID,
				Name:       month.Name,
				DaysCount:  month.DaysCount,
				ImageURL:   month.ImageURL,
				VideoURL:   month.VideoURL,
				LikesCount: len(month.Likes),
			},
		)
	}

	ctx.HTML(
		http.StatusOK,
		"month-list.html",
		gin.H{
			"months": items,
			"days":   daysString,
		},
	)
}

func (h *Handler) GetDepositMonth(ctx *gin.Context) {
	idString := strings.TrimSpace(ctx.Query("id"))

	var month repository.DepositMonth
	var err error

	if idString == "" {
		month, err = h.Repository.GetFirstPublishedDepositMonth()

		if err != nil {
			ctx.String(
				http.StatusNotFound,
				err.Error(),
			)
			return
		}
	} else {
		id, convertErr := strconv.Atoi(idString)

		if convertErr != nil {
			ctx.String(
				http.StatusBadRequest,
				"Некорректный ID",
			)
			return
		}

		if ctx.Query("next") == "true" {
			month, err =
				h.Repository.GetNextPublishedDepositMonth(id)
		} else {
			month, err =
				h.Repository.GetDepositMonth(id)
		}

		if err != nil {
			ctx.String(
				http.StatusNotFound,
				err.Error(),
			)
			return
		}
	}

	ctx.HTML(
		http.StatusOK,
		"month-details.html",
		gin.H{
			"month":      month,
			"likesCount": len(month.Likes),
		},
	)
}

func (h *Handler) GetDraftDepositMonth(ctx *gin.Context) {
	month, err :=
		h.Repository.GetDraftDepositMonth()

	if err != nil {
		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	ctx.HTML(
		http.StatusOK,
		"month-draft.html",
		gin.H{
			"month": month,
		},
	)
}
