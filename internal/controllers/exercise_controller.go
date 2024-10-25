package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/internal/models"
	"healthApi/internal/services"
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
	var exercise models.Exercise
	if err := ctx.ShouldBindJSON(&exercise); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := ctx.GetUint("userID")
	exercise.UserID = userID

	if err := c.service.CreateExercise(&exercise); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, exercise)
}

func (c *ExerciseController) GetUserExercises(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	log.Print("In GetUserExercises userId=", userID)
	exercises, err := c.service.GetUserExercises(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, exercises)
}

func (c *ExerciseController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var req models.Exercise
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
		ctx.JSON(http.StatusForbidden, "Not authorized to update this exercise record")
		return
	}

	// Update exercise fields
	exercise.Type = req.Type
	exercise.Duration = req.Duration
	exercise.Intensity = req.Intensity
	exercise.Date = req.Date
	exercise.Description = req.Description

	if err := c.service.Update(exercise); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := models.Exercise{

		Type:        exercise.Type,
		Duration:    exercise.Duration,
		Intensity:   exercise.Intensity,
		Date:        exercise.Date,
		Description: exercise.Description,
	}

	ctx.JSON(http.StatusOK, response)
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
