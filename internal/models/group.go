package models

type Group struct {
	Name        string `json:"name,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	GroupCode   string `json:"group_code,omitempty"`
	UID         string `json:"uid,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
}