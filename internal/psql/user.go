package psql

import (
	"fast/ent"
	"fast/ent/user"
	"fast/pkg/utils"
)

func mock_url() string {
	url_id := utils.Guid()
	photo_url := "https://re-up.ph/" + url_id + ".png"
	return photo_url
}

func NewUser(name string, email string, phone_number string, uid string, group_code string) string {

	group_id, err := GetGroupId(group_code)
	L.Fail("new-user", "get-group-id", group_id, err)

	photo_url := mock_url()

	user, err := pq.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPhoneNumber(phone_number).
		SetPhotoURL(photo_url).
		SetGroupCode(group_code).
		SetGroupID(group_id).
		SetIsActive(group_code != "NEO").
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

func CheckIfUserExists(uid string) (bool, *ent.User) {
	user, err := pq.User.Query().Where(user.UID(uid)).First(ctx)
	L.Fail("get-user", "by-uid", err)

	exists := false
	if user != nil {
		exists = true
		return exists, user
	}
	return exists, nil
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

/*
insert into groups
	(name,nickname,email,phone_number,group_code,photo_url,is_active,account_id,uid,address,group_id)
values
	('neo','neophyte','neo@fastinsure.com','+329154206969','NEO','https://re-up.ph/logo.png',true,'2aac7654-97be-499e-8e71-42179839365e','NEOxxxxxxxx000','re-up-hq','9f86a2c7-78a6-4b4f-8482-62077db90ee8');




	insert into groups (name,nickname,email,phone_number,group_code,photo_url,is_active,account_id,uid,address,id,address_id,create_time,update_time) values ('neo','neophyte','neo@fastinsure.com','+6329154206969','NEO','https://re-up.ph/logo.png',true,'39edf942-75e9-4bca-b71a-29c161be9b28','N7yCd3kCViMA0jD3eNuv5rqKxgy8--m','re-up-hq','9f86a2c7-78a6-4b4f-8482-62077db90ee8','re-up-hq-id',NOW(),NOW());
*/
