package psql

func CreateNewAccount(name string, email string, phone_number string, photo_url string, api_key string, uid string) (string, error) {

	account, err := pq.Account.
		Create().
		SetName(name).
		SetNickname(name).
		SetEmail(email).
		SetPhoneNumber(phone_number).
		SetPhotoURL(photo_url).
		SetAPIKey(api_key).
		SetUID(uid + "--m").
		SetAddressID("re-up-hq").
		SetIsActive(true).
		Save(ctx)

	L.Fail(r, "account-create", err)
	if err != nil {
		return "error", err
	}
	return account.UID, nil
}
