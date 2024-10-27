package controllers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services"
	utils2 "api-gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HydrationController struct {
	service services.HydrationService
}

func NewHydrationController(service services.HydrationService) *HydrationController {
	return &HydrationController{service: service}
}

func (c *HydrationController) Create(ctx *gin.Context) {
	var req models.HydrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse("Invalid request format : "+err.Error()))
		return
	}

	if errors := utils2.ValidateStruct(req); errors != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(errors[0].Error))
		return
	}

	userID := ctx.GetUint("userID")
	hydration := &models.Hydration{
		UserID: userID,
		Amount: req.Amount,
		Date:   time.Now(),
	}

	if err := c.service.Create(hydration); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to create hydration : "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils2.SuccessResponse("Hydration record created successfully", hydration.ToResponse()))
}
