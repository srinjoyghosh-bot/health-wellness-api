package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/api-gateway/internal/clients"
	"healthApi/api-gateway/internal/models"
	"healthApi/api-gateway/internal/services"
	utils2 "healthApi/api-gateway/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type ExerciseController struct {
	service services.ExerciseService
	client  clients.ExerciseClient
}

func NewExerciseController(service services.ExerciseService, client clients.ExerciseClient) *ExerciseController {
	return &ExerciseController{service: service, client: client}
}

func (c *ExerciseController) CreateExercise(ctx *gin.Context) {
	var req models.ExerciseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(err.Error()))
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

	_, err := c.client.CreateExercise(ctx.Request.Context(), userID.(uint), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse("Failed to create exercise : "+err.Error()))
		return
	}

	//var exercise models.Exercise = models.Exercise{
	//	UserID:      userID.(uint),
	//	Type:        req.Type,
	//	Duration:    req.Duration,
	//	Intensity:   req.Intensity,
	//	Date:        req.Date,
	//	Description: req.Description,
	//}
	//
	//if err := c.service.CreateExercise(&exercise); err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	ctx.JSON(http.StatusCreated, utils2.SuccessResponse("Exercise created", req))
}

func (c *ExerciseController) GetUserExercises(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	log.Print("In GetUserExercises userId=", userID)
	exercises, err := c.service.GetUserExercises(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse(err.Error()))
		return
	}

	var response []models.ExerciseResponse
	for _, exercise := range exercises {
		response = append(response, exercise.ToResponse())
	}

	ctx.JSON(http.StatusOK, utils2.SuccessResponse("Exercises found", response))
}

func (c *ExerciseController) Update(ctx *gin.Context) {
	id, err := utils2.ParseUint(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(err.Error()))
		return
	}

	var req models.ExerciseUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(err.Error()))
		return
	}

	if errors := utils2.ValidateStruct(req); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(errors[0].Error))
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
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse(
			"User not authenticated",
		))
		return
	}
	if exercise.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, utils2.ErrorResponse("Not authorized to update this exercise record"))
		return
	}

	// Update exercise fields
	exercise.UpdateFromRequest(req)

	if err := c.service.Update(exercise); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := exercise.ToResponse()

	ctx.JSON(http.StatusOK, utils2.SuccessResponse("Exercise updated", response))
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
