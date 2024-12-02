package api

import (
	"fast/ent"
	"fast/internal/models"
	"fast/internal/psql"

	"github.com/gofiber/fiber/v2"
)

func CreateGroup(c *fiber.Ctx) error {
	var v *models.Group
	if err := c.BodyParser(&v); err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}
	data := ent.Group{
		Name:        v.Name,
		Nickname:    v.Nickname,
		PhoneNumber: v.PhoneNumber,
		Email:       v.Email,
		PhotoURL:    &v.PhotoURL,
	}

	result := psql.CreateNewGroup(data.Name, data.Email, data.PhoneNumber, v.UID, v.GroupCode, data.AccountID, v.PhotoURL, true)

	return OkResponse(c, result, nil)
}

// if err != nil {
// 	return ErrResponse(c, BadRequest, err, JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
// }
