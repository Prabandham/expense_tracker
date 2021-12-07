package main

import (
	"log"
	"net/http"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/endpoints"
	"github.com/Prabandham/expense_tracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
	gin.ForceConsoleColor()
	router := gin.Default()
	router.Use(cors.Default())
	api := router.Group("/api/v1")
	// Unauthenticated routes
	api.POST("/login", endpoints.Login)
	api.POST("/register", endpoints.Register)

	// Authorized routes
	api.DELETE("/logout", TokenAuthMiddleware(redis), endpoints.Logout)
	api.GET("/credit_types", TokenAuthMiddleware(redis), endpoints.ListCreditTypes)
	api.POST("/credit_types", TokenAuthMiddleware(redis), endpoints.CreateCreditType)
	api.GET("/debit_types", TokenAuthMiddleware(redis), endpoints.ListDebitTypes)
	api.POST("/debit_types", TokenAuthMiddleware(redis), endpoints.CreateDebitType)
	api.GET("/credits", TokenAuthMiddleware(redis), endpoints.ListCredits)
	api.POST("/credits", TokenAuthMiddleware(redis), endpoints.CreateCredit)
	api.GET("/debits", TokenAuthMiddleware(redis), endpoints.ListDebits)
	api.POST("/debits", TokenAuthMiddleware(redis), endpoints.CreateDebit)
	api.GET("/accounts", TokenAuthMiddleware(redis), endpoints.ListAccounts)
	api.POST("/accounts", TokenAuthMiddleware(redis), endpoints.CreateAccount)

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
