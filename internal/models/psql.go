package models

import (
	"context"
	"fast/config"
	"fast/pkg/utils"
	"time"
)

var (
	pq  = config.LoadConfig().PQ
	ctx = context.Background()
	r   = "p_sql"
)

func PsqlHealth() interface{} {

	start := time.Now()

	err := pq.Ping()
	utils.WarnLog(r, "ping", "down", err)
	utils.Ok(r, "ping", "up")

	elapsed := time.Now().Sub(start) / time.Millisecond

	return map[string]interface{}{
		"sys":     "turso",
		"elapsed": elapsed,
		"unit":    "ms",
	}
}

func GetAccountWithAPIKey(apiKey string) (*Account, error) {

	var account Account
	query := "SELECT * FROM accounts WHERE api_key = $1"

	err := pq.QueryRow(query, apiKey).Scan(&account.UID, &account.Name, &account.Email, &account.Active, &account.CrTime, &account.Role, &account.APIKey)
	utils.NoRowsErrLog("api-key", "query-row", err)

	utils.OkLog(r, "is_active", account.Active, err)
	return &account, nil
}

func CheckAdminPrivileges(uid string) bool {
	var account Account
	query := "SELECT * FROM accounts WHERE uid = $1"

	matched := false
	err := pq.QueryRow(query, uid).Scan(&account.UID, &account.Name, &account.Email, &account.Active, &account.CrTime, &account.Role, &account.APIKey)
	utils.ErrLog(r, "uid", err)

	if uid == account.UID {
		matched = true
		utils.OkLog(r, "uid", "matched", err)
	}
	return matched
}

func QueryUsers() {
	rows, err := pq.Query("SELECT * FROM users")
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
