package service

import (
	"fast/internal/models"
	"fast/internal/psql"
)

func CreateNewAccount(a *models.Account) error {

	account_id, err := psql.CreateNewAccount(a.Name, a.Email, a.PhoneNumber, a.PhotoURL, a.APIKey, a.UID)

	L.Fail(r, "create-new account", err)
	if err != nil {
		return err
	}

	L.Good(r, "create-new account id:", account_id)
	return nil
}

func CreateNewGroup(g *models.Group) error {

	group_id, err := psql.CreateNewGroup(g.Name, g.Email, g.PhoneNumber, g.UID, g.GroupCode, g.AccountId, g.PhotoURL)

	L.Fail(r, "create-new group", err)
	if err != nil {
		return err
	}

	L.Good(r, "create-new group id:", group_id)
	return nil
}
