package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/internal/models"
	"healthApi/internal/services"
	"healthApi/internal/utils"
	"net/http"
	"time"
)

type SleepController struct {
	service services.SleepService
}

func NewSleepController(service services.SleepService) *SleepController {
	return &SleepController{service: service}
}

func (c *SleepController) Create(ctx *gin.Context) {
	var req models.SleepRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format : "+err.Error()))
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors[0].Error))
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	sleep := &models.Sleep{
		UserID:    userID.(uint),
		SleepTime: req.SleepTime,
		WakeTime:  req.WakeTime,
		Quality:   req.Quality,
		CreatedAt: time.Now(),
	}

	if err := c.service.Create(sleep); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create sleep record : "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Sleep record created successfully", sleep.ToResponse()))
}
