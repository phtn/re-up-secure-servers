package models

import "time"

type HCodeParams struct {
	KeyCode string `json:"key_code"`
	Code    string `json:"code"`
	GrpCode string `json:"grp"`
	Nonce   string `json:"nonce"`
	Sha     string `json:"sha"`
}

type HCodeVerification struct {
	Verified  bool           `json:"verified"`
	Expiry    *time.Duration `json:"expiry,omitempty"`
	Url       interface{}    `json:"url,omitempty"`
	GroupCode string         `json:"group_code,omitempty"`
}

type HCodeResponse struct {
	Code   string         `json:"code"`
	URL    string         `json:"url"`
	Expiry *time.Duration `json:"expiry"`
}
