package service

import (
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/shield"
)

func VerifyAgentCode(p *models.HCodeParams) *models.HCodeVerification {

	decrypt := shield.DecodeBase64(p.GrpCode)
	grp := shield.Decrypt(decrypt, p.KeyCode)
	L.Info("grp", "GrpCode", p.GrpCode, p.KeyCode)
	L.Info("grp", "decrypted", string(decrypt))

	v := rdb.RetrieveVal(p.Code + "--" + string(grp) + "--" + p.Nonce)

	if v.Exists {
		return &models.HCodeVerification{Verified: true, Expiry: &v.TTL, Url: &v.Value, GroupCode: string(grp)}
	}

	return &models.HCodeVerification{Verified: false}

}
