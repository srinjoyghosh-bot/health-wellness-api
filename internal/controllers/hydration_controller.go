package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/internal/models"
	"healthApi/internal/services"
	"healthApi/internal/utils"
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
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format : "+err.Error()))
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors[0].Error))
		return
	}

	userID := ctx.GetUint("userID")
	hydration := &models.Hydration{
		UserID: userID,
		Amount: req.Amount,
		Date:   time.Now(),
	}

	if err := c.service.Create(hydration); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create hydration : "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Hydration record created successfully", hydration.ToResponse()))
}
