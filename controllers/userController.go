package controllers

import (
	"go-jwt/db"
	"go-jwt/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context){
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		PayloadValidation(c, err)
		return
	}	

	if err := user.HasPassword(user.Password); err != nil {
		PayloadError(c, http.StatusInternalServerError, err.Error())
		return
	}

	record := db.Gorm.Create(&user)
	if record.Error != nil {
		PayloadError(c, http.StatusInternalServerError, record.Error.Error())
		return
	}
	user.Password = ""
	Payload(c, &user)
}

func GetUsers(c *gin.Context) {
	var data []models.User

    q := db.Gorm

	q = FilterBy(c, q, "id")
	q = SearchBy(c, q, "name")
	q = SearchBy(c, q, "username")
	q = FilterBy(c, q, "email")
	
	q.Omit("password").Find(&data)

	if len(data) <= 0 {
		Payload(c, nil)
	} else {
		Payload(c, &data)
	}
}

func UpdateUser(c *gin.Context) {
	var (
		data models.User
		dataUpdate models.UserUpdate
	)

	if err := c.ShouldBindJSON(&dataUpdate); err != nil {
		PayloadError(c, http.StatusBadRequest, err.Error())
		return
	}

    if err := db.Gorm.Where("id = ?", c.Param("id")).First(&data).Error; err != nil {
		Payload(c, nil)
    } else {
		data.Name = dataUpdate.Name

		record := db.Gorm.Save(&data)
		if record.Error != nil {
			PayloadError(c, http.StatusInternalServerError, record.Error.Error())
		} else {
			data.Password = ""
			Payload(c, &data)
		}
	}
}

func DeleteUser(c *gin.Context) {
	var data models.User

    if err := db.Gorm.Where("id = ?", c.Param("id")).First(&data).Error; err != nil {
		Payload(c, nil)
    } else {
		record := db.Gorm.Delete(&data)
		if record.Error != nil {
			PayloadError(c, http.StatusInternalServerError, record.Error.Error())
		} else {
			data.Password = ""
			Payload(c, &data)
		}
	}
}