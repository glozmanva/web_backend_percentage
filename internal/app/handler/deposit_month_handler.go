package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"deposit_month/internal/app/repository"
)

type Deposit_month_handler struct {
	Deposit_month_repository *repository.Deposit_month_repository
}

type Deposit_month_list_item struct {
	ID          int
	Name        string
	MonthNumber int
	DaysCount   int

	ImageURL string
	VideoURL string

	LikesCount int
}

func New_deposit_month_handler(
	deposit_month_repository *repository.Deposit_month_repository,
) *Deposit_month_handler {

	return &Deposit_month_handler{
		Deposit_month_repository: deposit_month_repository,
	}
}

func (h *Deposit_month_handler) Get_deposit_months(
	ctx *gin.Context,
) {

	min_days_count_string :=
		strings.TrimSpace(
			ctx.DefaultQuery(
				"min_days_count",
				"28",
			),
		)

	max_days_count_string :=
		strings.TrimSpace(
			ctx.DefaultQuery(
				"max_days_count",
				"31",
			),
		)

	min_days_count, err :=
		strconv.Atoi(
			min_days_count_string,
		)

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"Минимальное количество дней должно быть числом",
		)

		return
	}

	max_days_count, err :=
		strconv.Atoi(
			max_days_count_string,
		)

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"Максимальное количество дней должно быть числом",
		)

		return
	}

	if min_days_count < 28 {
		min_days_count = 28
	}

	if min_days_count > 31 {
		min_days_count = 31
	}

	if max_days_count < 28 {
		max_days_count = 28
	}

	if max_days_count > 31 {
		max_days_count = 31
	}

	if min_days_count > max_days_count {

		min_days_count,
			max_days_count =
			max_days_count,
			min_days_count
	}

	deposit_months :=
		h.Deposit_month_repository.
			Get_published_deposit_months(
				min_days_count,
				max_days_count,
			)

	deposit_month_items :=
		make(
			[]Deposit_month_list_item,
			0,
			len(deposit_months),
		)

	for _, deposit_month := range deposit_months {

		deposit_month_items =
			append(
				deposit_month_items,

				Deposit_month_list_item{
					ID: deposit_month.ID,

					Name: deposit_month.Name,

					MonthNumber: deposit_month.MonthNumber,

					DaysCount: deposit_month.DaysCount,

					ImageURL: deposit_month.ImageURL,

					VideoURL: deposit_month.VideoURL,

					LikesCount: len(
						deposit_month.Likes,
					),
				},
			)
	}

	min_days_position :=
		(min_days_count - 28) * 100 / 3

	max_days_position :=
		(max_days_count - 28) * 100 / 3

	ctx.HTML(
		http.StatusOK,
		"deposit_month_list.html",

		gin.H{
			"deposit_months": deposit_month_items,

			"min_days_count": min_days_count,

			"max_days_count": max_days_count,

			"min_days_position": min_days_position,

			"max_days_position": max_days_position,
		},
	)
}

func (h *Deposit_month_handler) Get_deposit_month(
	ctx *gin.Context,
) {

	deposit_month_id_string :=
		strings.TrimSpace(
			ctx.Query(
				"deposit_month_id",
			),
		)

	var deposit_month repository.Deposit_month
	var err error

	if deposit_month_id_string == "" {

		deposit_month, err =
			h.Deposit_month_repository.
				Get_first_published_deposit_month()

	} else {

		deposit_month_id,
			convert_error :=
			strconv.Atoi(
				deposit_month_id_string,
			)

		if convert_error != nil {

			ctx.String(
				http.StatusBadRequest,
				"Некорректный deposit_month_id",
			)

			return
		}

		if ctx.Query("next") == "true" {

			deposit_month, err =
				h.Deposit_month_repository.
					Get_next_published_deposit_month(
						deposit_month_id,
					)

		} else {

			deposit_month, err =
				h.Deposit_month_repository.
					Get_deposit_month(
						deposit_month_id,
					)
		}
	}

	if err != nil {

		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)

		return
	}

	likes_count :=
		len(
			deposit_month.Likes,
		)

	ctx.HTML(
		http.StatusOK,
		"deposit_month_feed.html",

		gin.H{
			"deposit_month": deposit_month,

			"likes_count": likes_count,

			"likes_count_after_like": likes_count + 1,
		},
	)
}

func (h *Deposit_month_handler) Get_draft_deposit_month(
	ctx *gin.Context,
) {

	deposit_month, err :=
		h.Deposit_month_repository.
			Get_draft_deposit_month()

	if err != nil {

		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"deposit_month_draft.html",

		gin.H{
			"deposit_month": deposit_month,
		},
	)
}
