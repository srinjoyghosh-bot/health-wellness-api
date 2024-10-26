package controllers

import (
	"github.com/gin-gonic/gin"
	"healthApi/internal/models"
	"healthApi/internal/services"
	"healthApi/internal/utils"
	"net/http"
)

type UserController struct {
	service services.UserService
	jwt     utils.JWTService
}

func NewUserController(service services.UserService, jwt utils.JWTService) *UserController {
	return &UserController{
		service: service,
		jwt:     jwt,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req models.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors[0].Error))
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

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("User registered", response))
}

func (c *UserController) Login(ctx *gin.Context) {
	var req models.UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	user, err := c.service.Authenticate(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid Credentials"))
		return
	}

	token, err := c.jwt.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	response := models.UserTokenResponse{TOKEN: token}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("User login", response))
}
