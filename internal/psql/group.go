package psql

import (
	"fast/ent/group"

	"github.com/google/uuid"
)

func CreateNewGroup(name string, email string, phone_number string, uid string, group_code string, account_id uuid.UUID, photo_url string) (string, error) {

	group, err := pq.Group.
		Create().
		SetName(name).
		SetNickname(name).
		SetEmail(email).
		SetPhoneNumber(phone_number).
		SetGroupCode(group_code).
		SetPhotoURL(photo_url).
		SetAccountID(account_id).
		SetUID(uid + "--m").
		SetAddress("re-up-hq").
		SetIsActive(true).
		Save(ctx)

	L.Fail(r, "group-create", err)

	if err != nil {
		return "error", err
	}

	return group.UID, nil
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
