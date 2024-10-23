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
	r   = "postgres"
	L   = utils.NewConsole()
)

func PsqlHealth() interface{} {

	start := time.Now()

	err := pq.Ping()
	L.Warn(r, "ping", "down", err)
	L.Good(r, "ping", "up", err)

	elapsed := time.Now().Sub(start) / time.Millisecond

	return map[string]interface{}{
		"sys":     "pq",
		"elapsed": elapsed,
		"unit":    "ms",
	}
}

func GetAccountWithAPIKey(apiKey string) (*Account, error) {

	var account Account
	query := "SELECT * FROM accounts WHERE api_key = $1"

	err := pq.QueryRow(query, apiKey).Scan(&account.UID, &account.Name, &account.Email, &account.Active, &account.CrTime, &account.Role, &account.APIKey)
	utils.NoRowsErrLog("accounts", "query-row", err)

	L.Good("accounts", "active", account.Active, err)
	return &account, nil
}

func CheckAdminPrivileges(uid string) bool {
	var account Account
	query := "SELECT * FROM accounts WHERE uid = $1"

	matched := false
	err := pq.QueryRow(query, uid).Scan(&account.UID, &account.Name, &account.Email, &account.Active, &account.CrTime, &account.Role, &account.APIKey)
	L.Fail(r, "query-by-uid", err)

	if uid == account.UID {
		matched = true
		L.Good(r, "matching-uid", "matched", err)
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
