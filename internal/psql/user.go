package psql

import (
	"fast/ent"
	"fast/ent/user"
)

func NewUser(name string, email string, phone_number string, uid string, group_code string) string {

	user, err := pq.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPhoneNumber(phone_number).
		SetGroupCode(group_code).
		SetIsActive(group_code != "").
		SetUID(uid).
		Save(ctx)

	L.Fail("users", "create", err)

	L.Good("users", "create", user.ID, err)
	return user.UID
}

func NewAccount(name string, email string, api_key string, uid string) string {

	account, err := pq.Account.
		Create().
		SetName(name).
		SetEmail(email).
		SetAPIKey(api_key).
		SetPhoneNumber("+639156984278").
		SetPhotoURL("https://fastinsure.com/images/re-up.png").
		SetUID(uid).
		Save(ctx)

	L.Fail("accounts", "create", err)
	L.Good("accounts", "create", account.UID, err)
	return account.UID
}

func CheckIfUserExists(uid string) bool {
	user, err := pq.User.Query().Where(user.UID(uid)).First(ctx)
	L.Fail("get-user", "by-uid", err)

	exists := false
	if user != nil {
		exists = true
	}
	return exists
}

func GetUserByUid(uid string) *ent.User {
	user, err := pq.User.Query().Where(user.UID(uid)).First(ctx)
	L.Fail("get-user", "by-uid", err)
	return user
}

func GetAllAccounts() {
	accounts, err := pq.Account.Query().All(ctx)
	L.Fail("all-accts", "query", err)
	L.Good("all-accts", "query", accounts, err)
	defer pq.Close()
}
