package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int
	Message string
	Data    interface{}
}

func HandleError(c *gin.Context, err interface{}) {
	respondWithError(err, c)
}

func HandleSuccess(c *gin.Context, data interface{}) {
	respondWithSuccess(data, c)
}

func respondWithError(err interface{}, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"data": err,
	})
}

func respondWithSuccess(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": &data,
	})
}
