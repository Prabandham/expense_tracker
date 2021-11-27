package endpoints

import (
	"net/http"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/gin-gonic/gin"
)

type LoginParams struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegesterParams struct {
	Email string `json:"email" binding:"required,email"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Authorized message"})
}

// Login endpoint will check credentials and generate a
// jwt - auth_token and refresh_token and send it back to the caller.
func Login(c *gin.Context) {
	var d LoginParams
	var user objects.User
	db := config.GetDatabaseConnection()
	redis := config.GetRedisConnection()

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	db.Connection.Where("email = ?", d.Email).First(&user)
	if !user.ValidatePassword(d.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Email or Password",
		})
		return
	}

	token, err := utils.CreateToken(user.ID.String()); if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := utils.CreateAuth(user.ID.String(), token, redis.Connection)
  if saveErr != nil {
     c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
  }

	c.JSON(200, gin.H{
		"access_token": token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

// Register endpoint will create a user
func Register(c *gin.Context) {
	var params RegesterParams 
	db := config.GetDatabaseConnection()

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := objects.User{Email: params.Email, Name: params.Name, Password: params.Password}
	result := db.Connection.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": result.Error,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Registered successfully",
	})
}

// Logout endpoint will delete the saved access and refresh tokens for a given user
func Logout(c *gin.Context) {
	redis := config.GetRedisConnection()
	au, err := utils.ExtractTokenMetadata(c.Request)
  if err != nil {
     c.JSON(http.StatusUnauthorized, "Unauthorized")
     return
  }
  deleted, delErr := utils.DeleteAuth(au.AccessUuid, redis.Connection)
  if delErr != nil || deleted == 0 { //if any goes wrong
     c.JSON(http.StatusUnauthorized, "Unauthorized")
     return
  }
  c.JSON(http.StatusOK, "Successfully logged out")
}
