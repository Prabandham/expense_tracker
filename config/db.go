package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	"github.com/Prabandham/expense_tracker/objects"
	"github.com/go-redis/redis/v7"
	_ "github.com/joho/godotenv/autoload"
)

type Db struct {
	Connection *gorm.DB
}

type Redis struct {
	Connection *redis.Client
}

type GormLogger struct{}

// Print - Log Formatter
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		log.WithFields(
			log.Fields{
				"module":        "gorm",
				"type":          "sql",
				"rows_returned": v[5],
				"src":           v[1],
				"values":        v[4],
				"duration":      v[2],
			},
		).Info(v[3])
	case "log":
		log.WithFields(log.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

var singleton *Db
var redisClient *Redis
var dbOnce sync.Once
var redisOnce sync.Once

// Get Env
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetDatabaseConnection() *Db {
	dbOnce.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			GetEnv("DB_HOST", "127.0.0.1"),
			GetEnv("DB_USER", ""),
			GetEnv("DB_NAME", ""),
			GetEnv("DB_PASSWORD", ""),
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
		dsn := GetEnv("REDIS_DSN", "localhost:6379")
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
	db.Connection.SetLogger(&GormLogger{})
	db.Connection.LogMode(true)
	formatter := new(log.JSONFormatter)
	log.SetFormatter(formatter)
	formatter.PrettyPrint = true
}

func (db *Db) MigrateModels() {
	db.Connection.AutoMigrate(
		&objects.User{},
		&objects.DebitType{},
		&objects.Debit{},
		&objects.Account{},
		&objects.CreditType{},
		&objects.Credit{},
	)
}
