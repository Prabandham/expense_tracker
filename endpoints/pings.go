package endpoints

import	"github.com/gin-gonic/gin"

func Ping(c *gin.Context) {
	HandleSuccess(c, gin.H{"data": "Pong !!!"})
}