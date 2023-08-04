package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ResponseData struct {
    Code 		int `json:"code" default:"200"`
    Message  	string `json:"message" default:"Success"`
	Data   		any `json:"data"`
}

type ResponseNoData struct {
    Code 		int `json:"code" default:"404"`
    Message  	string `json:"message" default:"Data Not Found"`
}

type ResponseValidation struct {
    Code 		int `json:"code" default:"400"`
    Message  	string `json:"message" default:"Bad Request"`
	Errors   	any `json:"errors"`
}

type ErrorMsg struct {
    Field string `json:"field"`
    Message   string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
    switch fe.Tag() {
        case "required":
            return "This field is required"
        case "lte":
            return "Should be less than " + fe.Param()
        case "gte":
            return "Should be greater than " + fe.Param()
    }
    return "Unknown error"
}


func Payload(c *gin.Context, T any){ 	
	if T==nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ResponseNoData{Code:404, Message:"Data Not Found"})
	} else {
		c.JSON(http.StatusOK, ResponseData{Code:200, Message:"Success", Data:T})
	}	
}

func PayloadError(c *gin.Context, code int, msg string) { 	
	c.AbortWithStatusJSON(code, ResponseNoData{Code:code, Message:msg})
}

func PayloadValidation(c *gin.Context, err error) { 		
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
		}				
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseValidation{Code:400, Message:"Validation Error", Errors:out})
		return
	}				

	c.AbortWithStatusJSON(http.StatusBadRequest, ResponseNoData{Code:400, Message:"Validation Error"})	
}

func FilterBy(c *gin.Context, q *gorm.DB, field string) (db *gorm.DB) { 
	db = q
	if c.Query(field) != "" {
		db = q.Where(field + " = ?", c.Query(field))
	}
	return
}

func SearchBy(c *gin.Context, q *gorm.DB, field string) (db *gorm.DB) { 
	db = q
	if c.Query(field) != "" {
		db = q.Where(field + " LIKE ?", "%"+c.Query(field)+"%")
	}
	return
}