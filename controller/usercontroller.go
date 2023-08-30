package controller

import (
	"github/golang-api/database"
	"github/golang-api/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserUpdate struct {
	Username string `json:"username" gorm:"unique" binding:"required"`
	Email    string `json:"email" gorm:"unique" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func RegisterUser(context *gin.Context) {
	var user model.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userID": user.ID, "username": user.Username, "email": user.Email})
}

func UpdateUser(context *gin.Context) {

	user, err := CurrentUser(context)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Where("id = ?", context.Param("id")).First(&user).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	var updateUser model.User

	if err := context.ShouldBindJSON(&updateUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := updateUser.HashPassword(updateUser.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if context.Param("id") != strconv.FormatUint(uint64(user.ID), 10) {
		context.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := database.Instance.Model(&user).Updates(model.User{Username: updateUser.Username, Email: updateUser.Email, Password: updateUser.Password}).Where("id = ?", user.ID).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"userID": user.ID, "username": user.Username, "email": user.Email})
}

func DeleteUser(context *gin.Context) {

	user, err := CurrentUser(context)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Where("id = ?", context.Param("id")).First(&user).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	if context.Param("id") != strconv.FormatUint(uint64(user.ID), 10) {
		context.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := database.Instance.Delete(&user).Where("id = ?", user.ID).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
