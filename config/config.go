package config

import (
	"context"
	"fast/pkg/utils"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

type Config struct {
	Addr string
	Fire *firebase.App
	Rdbs *redis.Client
}

func LoadConfig() Config {

	addr := os.Getenv("RE_UP_ADDR_PORT")

	fire := initialiazeFirebase()
	rdbs := initializeRedis()

	return Config{Addr: addr, Fire: fire, Rdbs: rdbs}
}

func initializeRedis() *redis.Client {

	rdbs := os.Getenv("RDB_SERV")
	host := os.Getenv("RDB_HOST")
	port := os.Getenv("RDB_PORT")
	pass := os.Getenv("RDB_PASS")

	opt, err := redis.ParseURL(rdbs + pass + host + port)
	utils.ErrLog("rdb", "config", err)

	rdb := redis.NewClient(opt)
	return rdb
}

func initialiazeFirebase() *firebase.App {

	cwd, err := os.Getwd()
	utils.ErrLog("fs", "cwd", err)

	pathToFile, exists := os.LookupEnv("SA_FILEPATH")
	if !exists {
		log.Fatal("SA_FILEPATH not set")
	}

	sa := filepath.Join(cwd, pathToFile)

	opt := option.WithCredentialsFile(sa)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	utils.ErrLog("init", "firebase", err)

	return app

}
