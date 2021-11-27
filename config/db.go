package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/go-redis/redis/v7"
)

type Db struct {
	Connection *gorm.DB
}

type Redis struct {
	Connection *redis.Client
}

var singleton *Db
var redisClient *Redis
var dbOnce sync.Once
var redisOnce sync.Once

func GetDatabaseConnection() *Db {
	dbOnce.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			utils.GetEnv("DB_HOST", "127.0.0.1"),
			utils.GetEnv("DB_USER", ""),
			utils.GetEnv("DB_NAME", ""),
			utils.GetEnv("DB_PASSWORD", ""),
		)
		db, err := gorm.Open("postgres", psqlInfo)
		if err != nil {
			panic("Could not connect to database")
		}
		singleton = &Db{Connection: db}
	})
	return singleton
}

func GetRedisConnection() *Redis {
	redisOnce.Do(func() {
		dsn := utils.GetEnv("REDIS_DSN", "localhost:6379")
		client := redis.NewClient(&redis.Options{
			Addr: dsn, //redis port
	  })
	  _, err := client.Ping().Result()
		if err != nil {
			panic("Could not connect to redis")
		}
		redisClient = &Redis{Connection: client}
	})
	return redisClient
}

func (db *Db) SetLogger() {
	db.Connection.LogMode(true)
	db.Connection.SetLogger(log.New(os.Stdout, "\r\n", 0))
}

func (db *Db) MigrateModels() {
	db.Connection.AutoMigrate(&objects.User{})
}
