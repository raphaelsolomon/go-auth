package handlers

import (
	"golang/internal/database"
	"golang/internal/jwt"
	"golang/internal/models"
	"golang/internal/password"
	"golang/internal/validation"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterHanders(ctx *gin.Context) {
	var inputDto struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		PhoneNo   string `json:"phone_no"`
	}

	if err := ctx.ShouldBindJSON((&inputDto)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate inputs
	if err := validation.ValidateEmail(inputDto.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validation.ValidatePassword(inputDto.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validation.ValidateName(inputDto.Firstname, "firstname"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validation.ValidateName(inputDto.Lastname, "lastname"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validation.ValidatePhoneNumber(inputDto.PhoneNo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPass := password.HashPassword(inputDto.Password)

	user := models.User{
		Email:           inputDto.Email,
		Firstname:       inputDto.Firstname,
		Lastname:        inputDto.Lastname,
		PhoneNo:         inputDto.PhoneNo,
		Status:          true,
		IsEmailverified: false,
		Password:        hashedPass,
		Uuid:            uuid.New().String(),
	}
	database := database.ConnectDatabase()
	if err := database.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func LoginHanders(ctx *gin.Context) {
	var inputDto struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON((&inputDto)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate inputs
	if err := validation.ValidateEmail(inputDto.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if inputDto.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	database := database.ConnectDatabase()

	var user models.User
	if database.Where("email = ?", inputDto.Email).First(&user).Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !password.ComparePassword(user.Password, inputDto.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := jwt.GenerateSignToken(user.Uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	responseBody := map[string]any{
		"access_token": accessToken,
		"expires_in":   time.Now().Add(time.Hour * 24),
		"token_type":   "Bearer",
		"user":         user.GetUser(),
	}

	ctx.JSON(http.StatusOK, responseBody)

}

func GetUserHandlers(ctx *gin.Context) {
	database := database.ConnectDatabase()

	var user models.User
	if database.Where("uuid = ?", ctx.GetString("user_uuid")).First(&user).Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user.GetUser())
}

func UpdateUserHandlers(ctx *gin.Context) {
	var inputDto struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		PhoneNo   string `json:"phone_no"`
	}

	if err := ctx.ShouldBindJSON((&inputDto)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate inputs (only validate if provided)
	if inputDto.Firstname != "" {
		if err := validation.ValidateName(inputDto.Firstname, "firstname"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if inputDto.Lastname != "" {
		if err := validation.ValidateName(inputDto.Lastname, "lastname"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if inputDto.PhoneNo != "" {
		if err := validation.ValidatePhoneNumber(inputDto.PhoneNo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	database := database.ConnectDatabase()
	if database.Where("uuid = ?", ctx.GetString("user_uuid")).Model(&models.User{}).Updates(inputDto).Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	var user models.User
	if database.Where("uuid = ?", ctx.GetString("user_uuid")).First(&user).Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No user found"})
		return
	}

	responseBody := map[string]any{
		"message": "User updated successfully",
		"user":    user.GetUser(),
	}

	ctx.JSON(http.StatusOK, responseBody)
}
