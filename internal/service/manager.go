package service

import (
	"fast/internal/models"
	"fast/internal/psql"
	"fast/internal/rdb"
	"fast/internal/shield"
	"fast/pkg/utils"
	"strings"
)

func NewAgentCode(v models.VerifyToken) *models.HCodeResponse {
	group_code := psql.GetGroupCode(v.UID)
	key_code, err := utils.GenerateCode()
	L.Fail("agent_code", "agent-code", err)
	L.Info("group_code", "psql", group_code)
	encryptedGrpCode := shield.EncodeBase64(shield.Encrypt([]byte(group_code), key_code))
	code := shield.NewKey(group_code)
	hcode := strings.Split(code, "--")
	encryptedUID := shield.Encrypt([]byte(v.UID), v.IDToken)
	encodedUID := shield.EncodeBase64(encryptedUID)

	dev_url := "http://localhost:3000"
	endpoint := "/hcode?code="
	url := dev_url + endpoint + hcode[0] + "&grp=" + encryptedGrpCode + "&nonce=" + hcode[2] + "&sha=" + encodedUID[:24]
	store_info := rdb.StoreVal(code, 48, url)
	L.Info("create  ", "agent-code", code, url, store_info.TTL, err)
	response := models.HCodeResponse{Code: key_code, URL: url, Expiry: &store_info.TTL}
	return &response
}
