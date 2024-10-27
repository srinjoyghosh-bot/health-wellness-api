package controllers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services"
	utils2 "api-gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	service services.UserService
	jwt     utils2.JWTService
}

func NewUserController(service services.UserService, jwt utils2.JWTService) *UserController {
	return &UserController{
		service: service,
		jwt:     jwt,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req models.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(err.Error()))
		return
	}
	if errors := utils2.ValidateStruct(req); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(errors[0].Error))
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := c.service.Create(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := c.jwt.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.UserTokenResponse{TOKEN: token}

	ctx.JSON(http.StatusCreated, utils2.SuccessResponse("User registered", response))
}

func (c *UserController) Login(ctx *gin.Context) {
	var req models.UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils2.ErrorResponse(err.Error()))
		return
	}

	user, err := c.service.Authenticate(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils2.ErrorResponse("Invalid Credentials"))
		return
	}

	token, err := c.jwt.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils2.ErrorResponse(err.Error()))
		return
	}

	response := models.UserTokenResponse{TOKEN: token}

	ctx.JSON(http.StatusOK, utils2.SuccessResponse("User login", response))
}
