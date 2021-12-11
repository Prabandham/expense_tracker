package endpoints

import (
	"errors"
	"net/http"

	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/gin-gonic/gin"
)

// Login endpoint will check credentials and generate a
// jwt - auth_token and refresh_token and send it back to the caller.
func Login(c *gin.Context) {
	var loginParams LoginParams
	c.BindJSON(&loginParams)

	var user objects.User
	if err := db.Connection.Where("email = ?", loginParams.Email).First(&user).Error; err != nil {
		HandleError(c, err)
		return
	}
	if !user.ValidatePassword(loginParams.Password) {
		HandleError(c, errors.New("invalid Username or Password"))
		return
	}

	token, err := utils.CreateAuth(user.ID.String(), redis.Connection)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleSuccess(c, gin.H{"access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}

// Register endpoint will create a user
func Register(c *gin.Context) {
	var params RegesterParams
	c.BindJSON(&params)

	user := objects.User{Email: params.Email, Name: params.Name, Password: params.Password}
	if err := db.Connection.Create(&user).Error; err != nil {
		HandleError(c, err)
		return
	}

	HandleSuccess(c, gin.H{"message": "Registered successfully"})
}

// Logout endpoint will delete the saved access and refresh tokens for a given user
func Logout(c *gin.Context) {
	au, err := utils.ExtractTokenMetadata(c.Request)
	HandleError(c, err)

	deleted, delErr := utils.DeleteAuth(au.AccessUuid, redis.Connection)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	HandleSuccess(c, "Successfully Logged out")
}
