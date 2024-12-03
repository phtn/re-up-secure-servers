package service

import (
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/shield"
)

func VerifyAgentCode(p *models.HCodeParams) models.HCodeVerification {

	decrypt := shield.DecodeBase64(p.Grp)
	grp := shield.Decrypt(decrypt, p.Code)
	L.Info(r, "hcode-grp-decrypt", string(grp))

	v := rdb.RetrieveVal(p.HKey + "--" + string(grp) + "--" + p.Nonce)
	L.Info(r, "rdb-get-value-exists", v.Value, v.Exists)

	if !v.Exists {
		return models.HCodeVerification{Verified: false}
	}

	return models.HCodeVerification{Verified: true, Expiry: &v.TTL, Url: &v.Value, GroupCode: string(grp)}

}
