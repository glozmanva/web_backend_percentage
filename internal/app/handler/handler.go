package handler

import (
	"deposit_month/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
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
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*.html")

	router.Static(
		"/static",
		"./resources",
	)
}
