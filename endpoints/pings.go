package endpoints

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	HandleSuccess(c, gin.H{"data": "Pong !!!"})
}

func BroadcastTest(c *gin.Context) {
	redisConnection := redis.Connection
	var statsParams StatsParams
	c.ShouldBindJSON(&statsParams)
	statsParams.RequestKey = CurrentUser(c).UserId
	payload, err := json.Marshal(statsParams)
	if err != nil {
		panic(err)
	}
	redisConnection.Publish("analytics", payload)
	HandleSuccess(c, gin.H{"data": "OK"})
}
