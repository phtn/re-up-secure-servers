package config

import (
	"context"
	"database/sql"
	"fast/pkg/utils"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"github.com/redis/go-redis/v9"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"google.golang.org/api/option"
)

type Config struct {
	Addr          string
	Fire          *firebase.App
	Rdbs          *redis.Client
	Db            *sql.DB
	ApiKey        string
	AllowedOrigin string
	JwtSecret     string
}

var (
	fire *firebase.App
	rdbs *redis.Client
	db   *sql.DB
	addr string
	akey string
	orig string
	jwts string
	turl string
	ttkn string
)

func init() {
	addr = os.Getenv("RE_UP_ADDR_PORT")
	akey = os.Getenv("RE_UP_API_KEY")
	orig = os.Getenv("RE_UP_ALLOWED_ORIGIN")
	jwts = os.Getenv("RE_UP_JWT_SECRET")
	turl = os.Getenv("TURSO_DATABASE_URL")
	ttkn = os.Getenv("TURSO_AUTH_TOKEN")
	fire = initialiazeFirebase()
	rdbs = initializeRedis()
	db = initialiazeDB()

	defer db.Close()

}

func LoadConfig() *Config {

	return &Config{Addr: addr, Fire: fire, Rdbs: rdbs, ApiKey: akey, AllowedOrigin: orig, JwtSecret: jwts, Db: db}
}

func initializeRedis() *redis.Client {

	rdba := os.Getenv("RDB_SERV")
	host := os.Getenv("RDB_HOST")
	port := os.Getenv("RDB_PORT")
	pass := os.Getenv("RDB_PASS")

	opt, err := redis.ParseURL(rdba + pass + host + port)
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

	fire, err := firebase.NewApp(context.Background(), nil, opt)
	utils.ErrLog("init", "firebase", err)

	return fire
}

func initialiazeDB() *sql.DB {
	url := turl + ttkn
	db, err := sql.Open("libsql", url)
	utils.ErrLog("db", "open", err)
	return db
}
