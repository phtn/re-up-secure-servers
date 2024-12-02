package api

import (
	"context"
	"fast/config"
	"fast/internal/rdb"
	"fast/pkg/utils"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

type Res struct {
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type DataResponse struct {
	Status  int
	Code    string
	Message string
	Err     error
	Data    interface{}
}

func Response(data interface{}, err error, message string) fiber.Map {
	response := Res{
		Data:    data,
		Error:   err,
		Message: message,
	}
	return fiber.Map{
		"data": response,
	}
}

var (
	f = utils.Raptor(" 𝐅 ", 0)
	z = config.LoadConfig().Zap
)

func destruct(appErr *AppError) (int, string, interface{}) {
	return appErr.Status, appErr.Message, appErr.Details
}

func ErrResponse(c *fiber.Ctx, appErr *AppError, err error) error {

	status, message, details := destruct(appErr)
	L.Fail(f, strconv.Itoa(status), details, err, message)
	if err != nil {
		L.Fail(f, strconv.Itoa(status), details, err, message)
	}
	return c.Status(status).JSON(appErr)

	// L.Good(f, strconv.Itoa(status), data.Details, data.Message)
	// return c.Status(status).JSON(data)
}

func OkResponse(c *fiber.Ctx, d interface{}, a interface{}) error {
	L.Good(f, strconv.Itoa(200), a)

	data := DataResponse{
		Status:  200,
		Code:    "ACCEPTED",
		Err:     nil,
		Message: "OK",
		Data:    d,
	}
	return c.Status(200).JSON(data)
	// L.Good(f, strconv.Itoa(status), data.Details, data.Message)
	// return c.Status(status).JSON(data)
}

func ValidateFields(u rdb.Tokens) *[]ValidationError {

	var errors []ValidationError

	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		isEmpty := value.IsZero()
		if isEmpty {
			errors = append(errors, ValidationError{
				Field:   field.Type.Name(),
				Message: fmt.Sprintf("Missing value in %v\n", field.Name),
			})
		}
	}

	if len(errors) > 0 {
		return &errors
	}

	return nil

}

func CookieHandler(c *fiber.Ctx, cookie string, data *DataResponse) error {

	if data.Err != nil {

		L.Fail(h, strconv.Itoa(data.Status), data.Err)
		return c.Status(data.Status).JSON(data)
	}

	expiresIn := time.Hour * 24 * 5

	c.Cookie(&fiber.Cookie{
		Name:     "fastinsure--session",
		Value:    cookie,
		Path:     "/",
		Expires:  time.Now().Add(expiresIn),
		Secure:   true,
		HTTPOnly: true,
	})

	L.Good(h, strconv.Itoa(data.Status), "OK")
	return c.Status(OK).JSON(data)
}

func validateToken(ctx context.Context, idToken string) (*auth.Token, error) {

	token, err := fire.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token %v\n", idToken)
	}

	L.Info(h, "validate-token", token.UID)

	if token.Expires < time.Now().Unix() {
		return nil, fmt.Errorf("token is expired %v\n", time.Until(time.Unix(token.Expires*int64(time.Second), token.Expires*int64(time.Nanosecond))))
	}

	return token, nil
}

// data := map[string]interface{}{
// 	"accepted": true,
// }
// return OkResponse(c, OK, DataResponse{
// 	Status:  OK,
// 	Code:    "Accepted",
// 	Message: "OK",
// 	Err:     nil,
// 	Data:    data,
// })
