package utils

import (
	"fmt"
	"encoding/json"

	"github.com/Prabandham/expense_tracker/config"
)

type Response struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
	RequestKey string `json:"request_key"`
}

func ListenForMessages() {
	redis := config.GetRedisConnection()
	subscriber := redis.Connection.Subscribe("analytics_response")

	for {
		msg, err := subscriber.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		response := Response{}

		if err := json.Unmarshal([]byte(msg.Payload), &response); err != nil {
			panic(err)
		}

		fmt.Println("Received message from " + msg.Channel + " channel.")
		fmt.Printf("%+v\n", response)

		payload, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}

		WebSocketHub.Broadcast <- Message{
			Data: payload,
			Room: response.RequestKey,
		}
	}
}