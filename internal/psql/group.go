package psql

import (
	"fast/ent/group"

	"github.com/google/uuid"
)

func CreateNewGroup(name string, email string, phone_number string, uid string, group_code string, account_id uuid.UUID, photo_url string, is_active bool) string {

	group, err := pq.Group.
		Create().
		SetName(name).
		SetNickname("").
		SetEmail(email).
		SetPhoneNumber(phone_number).
		SetGroupCode(group_code).
		SetPhotoURL(photo_url).
		SetIsActive(true).
		SetAccountID(account_id).
		SetUID(uid + "--m").
		SetAddress("re-up-hq").
		Save(ctx)

	L.Fail("groups", "create", err)

	L.Good("groups", "create", group.ID, err)
	return group.UID
}

func GetGroupCode(uid string) string {
	group, err := pq.Group.
		Query().Where(group.UID(uid + "--m")).First(ctx)
	if err != nil {
		L.Fail("groups", "get-code", err)
		return ""
	}
	return group.GroupCode
}
