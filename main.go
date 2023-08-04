package main

import (
	"fmt"
	"go-jwt/controllers"
	"go-jwt/db"
	"go-jwt/middlewares"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Initialize Database
	// database.Connect("root:root@tcp(localhost:3306)/db_jwt_go?parseTime=true")
	db.Connect(readEnvFileDb())
	db.Migrate()

	// Initialize Router
	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
	})
	
	api := router.Group("/api")
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"code": 200, "message": "Success", "os-env": runtime.GOOS +" - "+getEnvValue("APP_ENV", "development")})
    	})
		
		api.POST("/login", controllers.GenerateToken)
		api.POST("/users/register", controllers.RegisterUser)
		secured := api.Group("/v1").Use(middlewares.Auth())
		{
			secured.GET("/users", controllers.GetUsers)
			secured.PUT("/users/:id", controllers.UpdateUser)
			secured.DELETE("/users/:id", controllers.DeleteUser)
		}
	}	
	return router
}

func readEnvFileDb()(connectionString string){
	// reading environments variable using the viper
	appEnv := getEnvValue("APP_ENV", "development")
	// not set in our app.env
	appVersion := getEnvValue("APP_VERSION", "1")
	dbUser := getEnvValue("DB_USER", "root")
	dbPass := getEnvValue("DB_PASS", "")	
	dbIp := getEnvValue("DB_IP", "127.0.0.1")
	dbPort := getEnvValue("DB_PORT", "3306")
	dbName := getEnvValue("DB_NAME", "db_test")

	fmt.Printf(" ------%s-----\n", "Reading Environment variables")
	fmt.Printf(" %s = %s \n", "Application_Environment", appEnv)
	fmt.Printf(" %s = %s \n", "Application_Version", appVersion)
	fmt.Printf(" %s = %s \n", "Server_Listening_Address", dbIp)
	fmt.Printf(" %s = %s \n", "Server_Listening_Port", dbPort)
	fmt.Printf(" %s = %s \n", "Database_User_Name", dbUser)
	fmt.Printf(" %s = %s \n", "Database_User_Password", dbPass)	
	fmt.Printf(" %s = %s \n", "Database_Name", dbName)
	
	// db.Connect("root:@tcp(localhost:3306)/db_jwt_go?parseTime=true")
	fmt.Printf(" ConnectionString %s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbIp, dbPort, dbName)
	return dbUser + ":" + dbPass + "@tcp(" + dbIp + ":" + dbPort +")/" + dbName + "?parseTime=true"
}

func getEnvValue(key string, defaultValue string) string {
	// Get file path
	viper.SetConfigFile("app.env")
	//read configs and handle errors
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	value := viper.GetString(key)
	if value != "" {
		return value
	}
	return defaultValue
}
