package api

import (
	"fast/internal/rdb"
	"fast/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var (
	dev = utils.Dev("Δ", 0)
)

func DebugRedisStore(c *fiber.Ctx) error {

	var o rdb.TS

	if err := c.BodyParser(&o); err != nil {
		L.Fail(dev, "body-parser", err)
		return ErrResponse(c, ErrBadRequest, err)
	}

	L.Info(dev, o.Field, o.Value)
	key := fmt.Sprintf("%s", "dev-*-bright")
	// expiresIn := 5 * time.Minute
	store := rdb.DebugSet(key, o)

	L.Info(dev, "see what's up", store)

	return nil
}
