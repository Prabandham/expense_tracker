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
	respondWithError(err, c)
}

func HandleSuccess(c *gin.Context, data interface{}) {
	respondWithSuccess(data, c)
}

func respondWithError(err error, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"data": err.Error(),
	})
}

func respondWithSuccess(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": &data,
	})
}
