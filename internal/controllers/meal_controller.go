package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/internal/models"
	"healthApi/internal/services"
	"healthApi/internal/utils"
	"net/http"
)

type MealController struct {
	service services.MealService
}

func NewMealController(service services.MealService) *MealController {
	return &MealController{service: service}
}

func (c *MealController) Create(ctx *gin.Context) {
	var req models.MealRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format : "+err.Error()))
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors[0].Error))
		return
	}

	userID := ctx.GetUint("userID")
	meal := &models.Meal{
		UserID:      userID,
		Type:        req.Type,
		Description: req.Description,
		Date:        req.Date,
	}

	if err := c.service.Create(meal); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create Meal : "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Meal created succesfully", meal.ToResponse()))
}
