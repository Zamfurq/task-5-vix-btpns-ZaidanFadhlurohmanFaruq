package controller

import (
	"errors"
	"github/golang-api/database"
	"github/golang-api/helpers"
	"github/golang-api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PhotoUpdate struct {
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photoUrl" binding:"required"`
}

func CurrentUser(context *gin.Context) (model.User, error) {
	user_id, err := helpers.ExtractTokenID(context)

	if err != nil {
		return model.User{}, err
	}

	u, err := GetUserByID(user_id)

	if err != nil {
		return model.User{}, err
	}

	return u, nil
}

func GetUserByID(uid uint) (model.User, error) {

	var u model.User

	if err := database.Instance.Preload("Photos").Where("id=?", uid).Find(&u).Error; err != nil {
		return u, errors.New("User not found!")
	}

	return u, nil

}

func GetPhotos(context *gin.Context) {

	user, err := CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"photo": user.Photos})

}

func PostPhoto(context *gin.Context) {

	var photo model.Photo

	if err := context.ShouldBindJSON(&photo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := CurrentUser(context)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo.UserId = user.ID

	if err := database.Instance.Create(&photo).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, photo)
}

func UpdatePhoto(context *gin.Context) {

	var photo model.Photo

	user, err := CurrentUser(context)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Where("id = ?", context.Param("id")).First(&photo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Photo not found!"})
		return
	}

	var updatePhoto PhotoUpdate

	if err := context.ShouldBindJSON(&updatePhoto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if photo.UserId != user.ID {
		context.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := database.Instance.Model(&photo).Updates(model.Photo{Title: updatePhoto.Title, Caption: updatePhoto.Caption, PhotoUrl: updatePhoto.PhotoUrl}).Where("user_id = ?", user.ID).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, photo)
}

func DeletePhoto(context *gin.Context) {

	var photo model.Photo

	user, err := CurrentUser(context)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Where("id = ?", context.Param("id")).First(&photo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Photo not found!"})
		return
	}

	if photo.UserId != user.ID {
		context.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := database.Instance.Delete(&photo).Where("user_id = ?", user.ID).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Photo deleted"})
}
