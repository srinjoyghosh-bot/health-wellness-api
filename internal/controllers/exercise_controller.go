package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/internal/models"
	"healthApi/internal/services"
	"healthApi/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type ExerciseController struct {
	service services.ExerciseService
}

func NewExerciseController(service services.ExerciseService) *ExerciseController {
	return &ExerciseController{service: service}
}

func (c *ExerciseController) CreateExercise(ctx *gin.Context) {
	var req models.ExerciseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors[0].Error))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(
			"User not authenticated",
		))
		return
	}

	var exercise models.Exercise = models.Exercise{
		UserID:      userID.(uint),
		Type:        req.Type,
		Duration:    req.Duration,
		Intensity:   req.Intensity,
		Date:        req.Date,
		Description: req.Description,
	}

	if err := c.service.CreateExercise(&exercise); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Exercise created", req))
}

func (c *ExerciseController) GetUserExercises(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	log.Print("In GetUserExercises userId=", userID)
	exercises, err := c.service.GetUserExercises(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Exercises found", exercises))
}

func (c *ExerciseController) Update(ctx *gin.Context) {
	id, err := utils.ParseUint(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	var req models.ExerciseUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors[0].Error))
		return
	}

	// Verify ownership
	exercise, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	if exercise.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("Not authorized to update this exercise record"))
		return
	}

	// Update exercise fields
	exercise.UpdateFromRequest(req)

	if err := c.service.Update(exercise); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := models.ExerciseResponse{
		ID:          exercise.ID,
		Type:        exercise.Type,
		Duration:    exercise.Duration,
		Intensity:   exercise.Intensity,
		Date:        exercise.Date,
		Description: exercise.Description,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Exercise updated", response))
}

func (c *ExerciseController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Verify ownership
	exercise, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	userID, _ := ctx.Get("userID")
	if exercise.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, "\"Not authorized to delete this exercise record")
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Exercise record deleted successfully")
}
