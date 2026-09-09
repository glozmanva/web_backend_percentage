package api

import (
	"github.com/gin-gonic/gin"

	"web_backend_percentage/internal/app/handler"
	"web_backend_percentage/internal/app/repository"
)

func StartServer() error {
	repo := repository.NewRepository()
	h := handler.NewHandler(repo)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")


	router.Static("/static", "./resources")

	router.GET("/deposit-month", h.GetDepositMonth)
	router.GET("/deposit-month/draft", h.GetDraftDepositMonth)
	router.GET("/deposit-months", h.GetDepositMonths)

	return router.Run(":8080")
}
