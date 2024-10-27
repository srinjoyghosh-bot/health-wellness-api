package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/api-gateway/internal/models"
	"healthApi/api-gateway/internal/services"
	utils2 "healthApi/api-gateway/internal/utils"
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
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse("Invalid request format : "+err.Error()))
		return
	}

	if errors := utils2.ValidateStruct(req); errors != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(errors[0].Error))
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse(
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
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to create sleep record : "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils2.SuccessResponse("Sleep record created successfully", sleep.ToResponse()))
}
