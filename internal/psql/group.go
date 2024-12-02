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

	L.Fail(r, "group-create", err)

	L.Good(r, "group-create", group.ID, err)
	return group.UID
}

func GetGroupCode(uid string) string {
	group, err := pq.Group.
		Query().Where(group.UID(uid + "--m")).First(ctx)
	if err != nil {
		L.Fail(r, "group-get-code", err)
		return ""
	}
	return group.GroupCode
}

func GetGroupId(group_code string) (uuid.UUID, error) {
	group, err := pq.Group.
		Query().Where(group.GroupCode(group_code)).First(ctx)
	if err != nil {
		L.Fail(r, "get-group-id", group_code, err)
		return group.ID, err
	}
	return group.ID, nil
}
