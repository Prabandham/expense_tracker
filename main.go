package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/endpoints"
	"github.com/Prabandham/expense_tracker/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database, redis, setlogger and migrate models.
	redis := config.GetRedisConnection()
	db := config.GetDatabaseConnection()
	db.SetLogger()
	db.MigrateModels()

	// Start server and load routes
	gin.DisableConsoleColor()
	if config.GetEnv("GO_ENV", "development") == "production" {
		// Logging to a file.
		f, _ := os.Create("/home/ubuntu/apps/expense_tracker.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}

	router := gin.Default()
	corsConfig := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	router.Use(corsConfig)
	// Dummy router to check config / health of system
	go utils.WebSocketHub.Run()
	go utils.ListenForMessages()
	router.GET("/", endpoints.Ping)

	api := router.Group("/api/v1")
	// Unauthenticated routes
	api.POST("/login", endpoints.Login)
	api.POST("/register", endpoints.Register)

	// Authorized routes
	api.GET("/ws/:socketId", func(c *gin.Context) {
		socketId := c.Param("socketId")
		fmt.Println("Got socket ID as", socketId)
		utils.ServeWs(c.Writer, c.Request, socketId)
	})
	api.DELETE("/logout", TokenAuthMiddleware(redis), endpoints.Logout)
	api.GET("/credit_types", TokenAuthMiddleware(redis), endpoints.ListCreditTypes)
	api.POST("/credit_types", TokenAuthMiddleware(redis), endpoints.CreateCreditType)
	api.GET("/debit_types", TokenAuthMiddleware(redis), endpoints.ListDebitTypes)
	api.POST("/debit_types", TokenAuthMiddleware(redis), endpoints.CreateDebitType)
	api.DELETE("/debit_types/:id", TokenAuthMiddleware(redis), endpoints.DeleteDebitType)
	api.GET("/credits", TokenAuthMiddleware(redis), endpoints.ListCredits)
	api.POST("/credits", TokenAuthMiddleware(redis), endpoints.CreateCredit)
	api.DELETE("/credit_types/:id", TokenAuthMiddleware(redis), endpoints.DeleteCreditType)
	api.GET("/debits", TokenAuthMiddleware(redis), endpoints.ListDebits)
	api.POST("/debits", TokenAuthMiddleware(redis), endpoints.CreateDebit)
	api.GET("/accounts", TokenAuthMiddleware(redis), endpoints.ListAccounts)
	api.POST("/accounts", TokenAuthMiddleware(redis), endpoints.CreateAccount)
	api.GET("/accounts/:account_id/list_credits_and_debits", TokenAuthMiddleware(redis), endpoints.ListCreditsAndDebits)
	api.DELETE("/accounts/:id", TokenAuthMiddleware(redis), endpoints.DeleteAccount)
	api.GET("/grouped_credits", TokenAuthMiddleware(redis), endpoints.GroupedCredits)
	api.GET("/grouped_debits", TokenAuthMiddleware(redis), endpoints.GroupedDebits)

	// This endpoint will listen to calls made from the analytics engine and
	// then push that data to the user via websockets
	api.POST("/stats_callback", TokenAuthMiddleware(redis), endpoints.BroadcastTest)

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