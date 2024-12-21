package models

import "github.com/google/uuid"

type Group struct {
	Name        string    `json:"name,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	GroupCode   string    `json:"group_code,omitempty"`
	UID         string    `json:"uid,omitempty"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	AddressId   string    `json:"address_id,omitempty"`
	AccountId   uuid.UUID `json:"account_id,omitempty"`
	Active      string    `json:"is_active,omitempty"`
}
