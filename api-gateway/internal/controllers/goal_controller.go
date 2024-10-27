package controllers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services"
	utils2 "api-gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GoalController struct {
	service services.GoalService
}

func NewGoalController(service services.GoalService) *GoalController {
	return &GoalController{service: service}
}

func (c *GoalController) Create(ctx *gin.Context) {
	var req models.GoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse("Invalid request format : "+err.Error()))
		return
	}

	if errors := utils2.ValidateStruct(req); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(errors[0].Error))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	goal := &models.Goal{
		UserID:      userID.(uint),
		Type:        req.Type,
		Target:      req.Target,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Frequency:   req.Frequency,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := c.service.Create(goal); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to create goal : "+err.Error()))
		return
	}

	response := goal.ToResponse()

	ctx.JSON(http.StatusCreated, utils2.SuccessResponse("Goal created successfully", response))
}

func (c *GoalController) GetAll(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	goals, err := c.service.GetByUserID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to get all goals : "+err.Error()))
		return
	}

	var response []models.GoalResponse
	for _, goal := range goals {
		response = append(response, goal.ToResponse())
	}

	ctx.JSON(http.StatusOK, utils2.SuccessResponse("Goals retrieved successfully", response))
}

func (c *GoalController) GetActiveGoals(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	goals, err := c.service.GetActiveGoals(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to get active goals : "+err.Error()))
		return
	}

	var response []models.GoalResponse
	for _, goal := range goals {
		response = append(response, goal.ToResponse())
	}

	ctx.JSON(http.StatusOK, utils2.SuccessResponse("Active goals retrieved successfully", response))
}

func (c *GoalController) Update(ctx *gin.Context) {
	var req models.UpdateGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse("Invalid request format : "+err.Error()))
		return
	}

	goalID, e := utils2.ParseUint(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse("Invalid request format : "+e.Error()))
		return
	}
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse(
			"User not authenticated",
		))
		return
	}

	goal, err := c.service.GetByID(goalID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils2.ErrorResponse("Failed to get goal : "+err.Error()))
		return
	}

	if goal.UserID != userID {
		ctx.JSON(http.StatusForbidden, utils2.ErrorResponse("Not authorised to update this goal"))
		return
	}

	goal.UpdateFromRequest(req)

	if err := c.service.Update(goal); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to update goal : "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils2.SuccessResponse("Goal updated successfully", goal.ToResponse()))
}

func (c *GoalController) Delete(ctx *gin.Context) {
	goalID, e := utils2.ParseUint(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse("Invalid request format : "+e.Error()))
		return
	}
	userID := ctx.GetUint("userID")

	goal, err := c.service.GetByID(goalID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils2.ErrorResponse("Failed to get goal : "+err.Error()))
		return
	}

	if goal.UserID != userID {
		ctx.JSON(http.StatusForbidden, utils2.ErrorResponse("Not authorised to delete this goal"))
		return
	}

	if err := c.service.Delete(goalID); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to delete goal : "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils2.SuccessResponse("Goal deleted successfully", nil))
}
