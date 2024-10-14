package models

import (
	"context"
	"fast/config"
	"fast/pkg/utils"
	"time"
)

var (
	db  = config.LoadConfig().Db
	r   = "turso"
	ctx = context.Background()
)

func Ping() {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	utils.ErrLog(r, "ctx-down", err)
	utils.Ok(r, "ctx-status", "up")

	err = db.Ping()
	utils.ErrLog(r, "ping-down", err)
	utils.Ok(r, "ping", "pong")
}

func GetAccountAPIKey(apiKey string) (*Account, error) {
	var account Account
	query := "SELECT id, name, email FROM users WHERE api_key = ?"

	err := db.QueryRow(query, apiKey).Scan(&account.UID, &account.Name, &account.Email)
	utils.ErrLog(r, "api-key", err)
	utils.OkLog(r, "api-key", "matched", err)

	return &account, nil
}

func QueryUsers() {
	rows, err := db.Query("SELECT * FROM users")
	utils.ErrLog(r, "query", err)
	defer rows.Close()

	var users []Account

	for rows.Next() {
		var user Account

		err := rows.Scan(&user.UID, &user.Name)
		utils.ErrLog(r, "rows", err)

		users = append(users, user)

		utils.Ok(user.UID, user.Name, "")
	}

	utils.ErrLog(r, "users", err)

}
