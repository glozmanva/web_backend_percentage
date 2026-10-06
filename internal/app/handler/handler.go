package handler

import (
	"deposit_month/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(
	r *repository.Repository,
) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(
	router *gin.Engine,
) {
	// Старые маршруты интерфейса ЛР2
	router.GET(
		"/deposit_month",
		h.GetDepositMonth,
	)

	router.GET(
		"/deposit_month/draft",
		h.GetDraftDepositMonth,
	)

	router.GET(
		"/deposit_month/list",
		h.GetDepositMonths,
	)

	router.POST(
		"/deposit_month/draft",
		h.CreateDraftDepositMonth,
	)

	router.POST(
		"/deposit_month/publish",
		h.PublishDepositMonth,
	)

	router.POST(
		"/deposit_month/delete",
		h.DeleteDepositMonth,
	)

	// REST API ЛР3
	api := router.Group("/api")

	depositMonths := api.Group("/deposit-months")

	depositMonths.GET(
		"",
		h.GetDepositMonthsAPI,
	)

	depositMonths.GET(
		"/feed",
		h.GetDepositMonthFeedAPI,
	)

	depositMonths.GET(
		"/draft",
		h.GetDraftDepositMonthAPI,
	)

	depositMonths.POST(
		"",
		h.CreateDepositMonthAPI,
	)

	depositMonths.PUT(
		"/draft",
		h.PublishDepositMonthAPI,
	)

	depositMonths.DELETE(
		"/:id",
		h.DeleteDepositMonthAPI,
	)

	depositMonths.POST(
		"/:id/like",
		h.SetDepositMonthLikeAPI,
	)

	users := api.Group("/users")

	users.POST(
		"",
		h.RegisterUserAPI,
	)

	users.POST(
		"/authentication",
		h.AuthenticateUserAPI,
	)

	users.POST(
		"/deauthentication",
		h.DeauthenticateUserAPI,
	)
}

func (h *Handler) RegisterStatic(
	router *gin.Engine,
) {
	router.LoadHTMLGlob(
		"templates/*.html",
	)

	router.Static(
		"/static",
		"./resources",
	)
}
