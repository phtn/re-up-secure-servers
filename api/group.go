package api

import (
	"fast/ent"
	"fast/internal/models"
	"fast/internal/psql"
	"fast/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateGroup(c *fiber.Ctx) error {
	var v *models.Group
	if err := c.BodyParser(&v); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	data := ent.Group{
		Name:        v.Name,
		Nickname:    v.Nickname,
		PhoneNumber: v.PhoneNumber,
		Email:       v.Email,
		PhotoURL:    &v.PhotoURL,
	}

	result := psql.CreateNewGroup(data.Name, data.Email, data.PhoneNumber, v.UID, v.GroupCode, data.AccountID, v.PhotoURL, true)

	return utils.FiberResponse(c, utils.OK, nil, utils.JsonData{Data: map[string]interface{}{"group_uid": result, "status": "success"}})
}

// if err != nil {
// 	return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
// }
