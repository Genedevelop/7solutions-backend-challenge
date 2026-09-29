package middleware

import (
	"encoding/json"
	"errors"
	"io"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func Bind(c echov4.Context, dst any) error {
	// Step 1: decode exactly one JSON value from the body
	dec := json.NewDecoder(c.Request().Body)
	if err := dec.Decode(dst); err != nil {
		return errs.ErrInvalidJSON
	}

	// Step 2: anything after that value, like {"name":"x"} garbage, makes the whole body invalid
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return errs.ErrInvalidJSON
	}
	return nil
}
