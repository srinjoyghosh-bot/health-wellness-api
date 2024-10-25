package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/internal/models"
	"healthApi/internal/services"
	"log"
	"net/http"
)

type ExerciseController struct {
	service *services.ExerciseService
}

func NewExerciseController(service *services.ExerciseService) *ExerciseController {
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
