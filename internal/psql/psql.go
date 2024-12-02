package psql

import (
	"context"
	"fast/config"
	"fast/ent/account"
	"fast/pkg/utils"
	"time"

	_ "github.com/lib/pq"
)

var (
	pq  = config.LoadConfig().Pq
	ctx = context.Background()
	L   = utils.NewConsole()
	r   = utils.Sky(" 𝐏 ", 0)
)

func PsqlHealth() interface{} {

	start := time.Now()

	tx, err := pq.Tx(ctx)
	L.Warn(r, "ping", "down", err)
	L.Good(r, "ping", "up", tx, err)

	elapsed := time.Now().Sub(start) / time.Millisecond

	return map[string]interface{}{
		"sys":     "pq",
		"elapsed": elapsed,
		"unit":    "ms",
	}
}

func CheckAPIKey(api_key string) (bool, error) {

	account, err := pq.
		Account.
		Query().Where(
		account.APIKey(api_key),
		account.IsActive(true)).First(ctx)
	if err != nil {
		return false, err
	}

	L.Fail(r, "mother-account verified", err)
	return account.IsActive, nil
}

func CheckAdminPrivileges(uid string) bool {

	matched := false

	account, err := pq.Account.Query().Where(account.UID(uid)).Only(ctx)
	L.Fail(r, "query-by-uid", err)

	if uid == account.UID {
		matched = true
		L.Good(r, "matching-uid", "matched", err)
	}
	return matched
}
