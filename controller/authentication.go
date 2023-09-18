package controller

import (
	"github/be/common"
	"github/be/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONResult struct {
	Jwt string `json:"jwt"`
}

// @BasePath /api/v1

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Param input body model.AuthenticationInput true "User credentials"
// @Success 201 {object} model.User
// @Router /auth/register [post]
func Register(ctx *gin.Context) {
	var input model.AuthenticationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user := model.User{
		Username: input.Username,
		Password: input.Password,
	}

	savedUser, err := user.Save()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

// Login handles user login
// @Summary Login a user
// @Description Login a user
// @Accept json
// @Produce json
// @Param input body model.AuthenticationInput true "User credentials"
// @Success 200 {object} JSONResult
// @Router /auth/login [post]
func Login(ctx *gin.Context) {
	var input model.AuthenticationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	user, err := model.FindUserByUsername(input.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
		return
	}

	err = user.ValidatePassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
		return
	}

	jwt, err := common.GenerateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
