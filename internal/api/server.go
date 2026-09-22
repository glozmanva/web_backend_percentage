package api

import (
	"github.com/gin-gonic/gin"

	"deposit_month/internal/app/handler"
	"deposit_month/internal/app/repository"
)

func StartServer() error {
	deposit_month_repository := repository.New_deposit_month_repository()

	deposit_month_handler :=
		handler.New_deposit_month_handler(deposit_month_repository)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.Static(
		"/static",
		"./resources",
	)

	router.GET(
		"/deposit_month",
		deposit_month_handler.Get_deposit_month,
	)

	router.GET(
		"/deposit_month/draft",
		deposit_month_handler.Get_draft_deposit_month,
	)

	router.GET(
		"/deposit_month/list",
		deposit_month_handler.Get_deposit_months,
	)

	return router.Run(":8080")
}
