package endpoints

import (
	"errors"
	"net/http"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/gin-gonic/gin"
)

// Login endpoint will check credentials and generate a
// jwt - auth_token and refresh_token and send it back to the caller.
func Login(c *gin.Context) {
	var d LoginParams
	var user objects.User
	db := config.GetDatabaseConnection()
	redis := config.GetRedisConnection()

	HandleError(c, c.ShouldBindJSON(&d))

	db.Connection.Where("email = ?", d.Email).First(&user)
	if !user.ValidatePassword(d.Password) {
		HandleError(c, errors.New("invalid Username or Password"))
	}

	token, err := utils.CreateToken(user.ID.String())
	HandleError(c, err)

	saveErr := utils.CreateAuth(user.ID.String(), token, redis.Connection)
	if saveErr != nil {
		HandleError(c, errors.New(saveErr.Error()))
	}

	HandleSuccess(c, gin.H{"access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}

// Register endpoint will create a user
func Register(c *gin.Context) {
	var params RegesterParams
	db := config.GetDatabaseConnection()

	HandleError(c, c.ShouldBindJSON(&params))

	user := objects.User{Email: params.Email, Name: params.Name, Password: params.Password}
	result := db.Connection.Create(&user)

	HandleError(c, result.Error)

	HandleSuccess(c, "Registered successfully")
}

// Logout endpoint will delete the saved access and refresh tokens for a given user
func Logout(c *gin.Context) {
	redis := config.GetRedisConnection()
	au, err := utils.ExtractTokenMetadata(c.Request)
	HandleError(c, err)

	deleted, delErr := utils.DeleteAuth(au.AccessUuid, redis.Connection)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	HandleSuccess(c, "Successfully Logged out")
}
