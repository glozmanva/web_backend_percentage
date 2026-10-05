package handler

import (
	"errors"
	"net/http"
	"time"

	"deposit_month/internal/app/ds"
	"deposit_month/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUserAPI(
	ctx *gin.Context,
) {
	var request struct {
		FullName  string `json:"full_name" binding:"required"`
		BirthDate string `json:"birth_date" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
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

	birthDate, err := time.Parse(
		"2006-01-02",
		request.BirthDate,
	)

	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "birth_date должен иметь формат YYYY-MM-DD",
			},
		)
		return
	}

	user := ds.User{
		FullName:  request.FullName,
		BirthDate: birthDate,
		Email:     request.Email,
		Password:  request.Password,
	}

	err = h.Repository.CreateUser(&user)

	if errors.Is(
		err,
		repository.ErrEmailExists,
	) {
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"error": err.Error(),
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
		http.StatusCreated,
		gin.H{
			"id":         user.ID,
			"full_name":  user.FullName,
			"birth_date": request.BirthDate,
			"email":      user.Email,
		},
	)
}

func (h *Handler) AuthenticateUserAPI(
	ctx *gin.Context,
) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "заглушка аутентификации для ЛР4",
		},
	)
}

func (h *Handler) DeauthenticateUserAPI(
	ctx *gin.Context,
) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "заглушка деавторизации для ЛР4",
		},
	)
}
