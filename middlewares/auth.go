package middlewares

import (
	"go-jwt/auth"
	"go-jwt/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc{
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {			
			controllers.PayloadError(c, http.StatusUnauthorized, "Request does not contain an access token")
			return
		}
		err:= auth.ValidateToken(tokenString)
		if err != nil {
			controllers.PayloadError(c, http.StatusUnauthorized, err.Error())			
			return
		}
		c.Next()
	}
}
