package config

import (
	"context"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"github.com/go-redis/redis/v8"
	"google.golang.org/api/option"
)

type Config struct {
	ServerAddress string
	App           *firebase.App
	Rdb           *redis.Client
}

func LoadConfig() Config {

	// addr, exists := os.LookupEnv("ADDR")

	// if !exists {
	// 	addr = ":1981"
	// }

	app := initialiazeFirebase()
	rdb := initializeRedis()

	return Config{ServerAddress: ":1981", App: app, Rdb: rdb}
}

func initializeRedis() *redis.Client {
	adr, exist := os.LookupEnv("REDIS_PORT")

	if !exist {
		adr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: adr,
	})
	return rdb
}

func initialiazeFirebase() *firebase.App {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory", err)
	}

	pathToFile, exists := os.LookupEnv("SA_FILEPATH")
	if !exists {
		log.Fatal("SA_FILEPATH not set")
	}

	sa := filepath.Join(cwd, pathToFile)

	opt := option.WithCredentialsFile(sa)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("Error initializing app", err)
	}
	return app

}
