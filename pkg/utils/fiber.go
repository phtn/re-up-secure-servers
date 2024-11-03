package utils

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Response(data interface{}, err error, message string) fiber.Map {
	response := JsonData{
		Data:    data,
		Error:   err,
		Message: message,
	}
	return fiber.Map{
		"data": response,
	}
}

func FiberResponse(c *fiber.Ctx, status int, err error, data JsonData) error {
	if err != nil {
		L.Fail(strconv.Itoa(status), data.Message, data.Data, err)
		return c.Status(status).JSON(Response(data, err, data.Message))
	}
	L.Good(strconv.Itoa(status), data.Message, data.Data)
	return c.Status(status).JSON(data)
}

func FiberCookie(c *fiber.Ctx, cookie string, status int, err error, data JsonData) error {
	if err != nil {
		L.Fail(strconv.Itoa(status), data.Message, data.Data, err)
		return c.Status(status).JSON(Response(data, err, data.Message))
	}
	L.Good(strconv.Itoa(status), data.Message, data.Data)

	expires := 287 * time.Hour

	c.Cookie(&fiber.Cookie{
		Name:     "fastinsure--cookie-session",
		Value:    cookie,
		Path:     "/authed",
		Expires:  time.Now().Add(expires),
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(status).JSON(data)
}
