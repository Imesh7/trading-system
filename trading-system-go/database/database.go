package database

import (
	//"c/models/balance"
	//"c/models/order"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	DataBase *gorm.DB
}

type RedisDBInstance struct {
	Client *redis.Client
}

var DB DBInstance
var RedisDB RedisDBInstance

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	fmt.Print(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Printf("Cannot connect database %e", err)
		os.Exit(1)
	}

	DB = DBInstance{
		DataBase: db,
	}
}

func ConnectToRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "1234",
		DB:       0, // use default DB
	})

	RedisDB = RedisDBInstance{
		Client: client,
	}
}
