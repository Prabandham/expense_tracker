package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status int
	Message string
	Data interface{}
}

func HandleError(c *gin.Context, err error) {
	if err != nil {
		respondWithError(http.StatusUnprocessableEntity, err, c)
		return
	}
}

func HandleSuccess(c *gin.Context, data interface{}) {
	respondWithSuccess(data, c)
}

func respondWithError(status int, err error, c *gin.Context) {
	c.JSON(http.StatusUnprocessableEntity, err)
}

func respondWithSuccess(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": &data,
	})
}
