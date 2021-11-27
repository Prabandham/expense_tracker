package main

import (
	"log"
	"net/http"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/Prabandham/expense_tracker/endpoints"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := env.Load()
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	// Initialize database, redis, setlogger and migrate models.
	redis := config.GetRedisConnection()
	db := config.GetDatabaseConnection()
	db.SetLogger()
	db.MigrateModels()

	// Start server and load routes
	router := gin.Default()
	api := router.Group("/api/v1")
	// Unauthenticated routes
	api.POST("/login", endpoints.Login)
	api.POST("/register", endpoints.Register)

	api.GET("/hello", TokenAuthMiddleware(redis), endpoints.Hello)
	api.DELETE("/logout", TokenAuthMiddleware(redis), endpoints.Logout)

	log.Fatal(router.Run(":3000"))
}

func TokenAuthMiddleware(redis *config.Redis) gin.HandlerFunc {
  return func(c *gin.Context) {
     err := utils.TokenValid(c.Request)
     if err != nil {
        c.JSON(http.StatusUnauthorized, "Unauthorized")
        c.Abort()
        return
     }
		 tokenAuth, err := utils.ExtractTokenMetadata(c.Request)
		 if err != nil {
				c.JSON(http.StatusUnauthorized, "unauthorized")
				return
		 }
		 _, err = utils.FetchAuth(tokenAuth, redis.Connection)
		 if err != nil {
		 		c.JSON(http.StatusUnauthorized, "Unauthorized")
        c.Abort()
				return
		}
     c.Next()
  }
}
