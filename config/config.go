package config

import (
	"context"
	"database/sql"
	"fast/ent"
	"fast/ent/migrate"
	"fast/pkg/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"

	dialect "entgo.io/ent/dialect"
	esql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/redis/go-redis/v9"

	// _ "github.com/tursodatabase/libsql-client-go/libsql"
	"google.golang.org/api/option"

	"go.uber.org/zap"
)

type Firebase struct {
	AuthClient *auth.Client
}

type Config struct {
	Addr          string
	Fire          *Firebase
	Rdbs          *redis.Client
	Pq            *ent.Client
	Zap           *zap.Logger
	ApiKey        string
	AllowedOrigin string
	JwtSecret     string
}

var (
	fire *Firebase
	rdbs *redis.Client
	pq   *ent.Client
	z    *zap.Logger
	addr string
	akey string
	orig string
	jwts string
	turl string
	ttkn string
	pdsn string
	L    = utils.NewConsole()
)

func init() {
	addr = os.Getenv("RE_UP_ADDR_PORT")
	akey = os.Getenv("RE_UP_API_KEY")
	orig = os.Getenv("RE_UP_ALLOWED_ORIGIN")
	jwts = os.Getenv("RE_UP_JWT_SECRET")
	turl = os.Getenv("TURSO_DATABASE_URL")
	ttkn = os.Getenv("TURSO_AUTH_TOKEN")
	pdsn = os.Getenv("SB_DSN")
	fire = initFirebase()
	rdbs = initRedis()
	pq = initPostgres()

}

func LoadConfig() *Config {

	return &Config{Addr: addr, Fire: fire, Rdbs: rdbs, ApiKey: akey, AllowedOrigin: orig, JwtSecret: jwts, Pq: pq, Zap: z}
}

func initRedis() *redis.Client {

	rdba := os.Getenv("RDB_SERV")
	host := os.Getenv("RDB_HOST")
	port := os.Getenv("RDB_PORT")
	pass := os.Getenv("RDB_PASS")

	opt, err := redis.ParseURL(rdba + pass + host + port)
	L.Fail("rdb", "config", err)

	rdb := redis.NewClient(opt)
	return rdb
}

func initFirebase() *Firebase {

	cwd, err := os.Getwd()
	L.Fail("fs", "cwd", err)

	pathToFile, exists := os.LookupEnv("SA_FILEPATH")
	if !exists {
		log.Fatal("SA_FILEPATH not set")
	}

	sa := filepath.Join(cwd, pathToFile)

	opt := option.WithCredentialsFile(sa)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	L.Fail("init", "firebase-new-app", err)

	client, err := app.Auth(context.Background())
	L.Fail("init", "firebase-auth", err)

	return &Firebase{AuthClient: client}
}

// func initialiazeDB() *sql.DB {
// 	// dataSourceName := turl + ttkn
// 	dataSourceName := "file:./local.db"
// 	db, err := sql.Open("libsql", dataSourceName)
// 	L.Fail("db", "open", err)
// 	return db
// }

func initPostgres() *ent.Client {
	dataSourceName := pdsn
	// dataSourceName := "postgres://xpriori:phtn458@localhost:5432/dpqb?sslmode=disable"
	db, err := sql.Open("postgres", dataSourceName)
	L.Fail("pq", "open", err)

	driver := dialect.Postgres
	client := ent.NewClient(ent.Driver(esql.OpenDB(driver, db)))
	ctx := context.Background()

	options := []schema.MigrateOption{
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
		migrate.WithForeignKeys(true),
		migrate.WithGlobalUniqueID(true),
	}

	client.Schema.Create(ctx, options...)

	return client
}

func initZap() *zap.Logger {
	raw_json := []byte(`{
		  "level": "debug",
		  "encoding": "json",
		  "outputPaths": ["stdout", "/tmp/logs"],
		  "errorOutputPaths": ["stderr"],
		  "initialFields": {"foo": "bar"},
		  "encoderConfig": {
		    "messageKey": "message",
		    "levelKey": "level",
		    "levelEncoder": "lowercase"
		  }
		}`)

	var cfg zap.Config
	if err := json.Unmarshal(raw_json, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	return logger

}
