package middleware

import (
	"errors"
	"net/http"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/logger"
)

func ErrorHandler() echov4.HTTPErrorHandler {
	return func(err error, c echov4.Context) {
		// Step 1: the logger middleware and Echo can both call this, write only once
		if c.Response().Committed {
			return
		}

		// Step 2: turn Echo's own errors (unknown route, bad method, body too large) into central errors
		var herr *echov4.HTTPError
		if errors.As(err, &herr) {
			err = fromEchoError(herr)
		}
		e := errs.From(err)

		// Step 3: 5xx means a bug or an outage, log the real cause since the client only sees a generic message
		status := e.HTTPStatus()
		if status >= http.StatusInternalServerError {
			logger.Error("request failed", e,
				"request_id", c.Response().Header().Get(echov4.HeaderXRequestID),
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
			)
		}
		_ = c.JSON(status, e.Response())
	}
}

func fromEchoError(herr *echov4.HTTPError) error {
	switch herr.Code {
	case http.StatusNotFound:
		return errs.ErrRouteNotFound
	case http.StatusMethodNotAllowed:
		return errs.ErrMethodNotAllowed
	case http.StatusRequestEntityTooLarge:
		return errs.ErrBodyTooLarge
	}
	if herr.Code >= http.StatusInternalServerError {
		return errs.Internal(herr, "echo")
	}
	msg, ok := herr.Message.(string)
	if !ok {
		msg = http.StatusText(herr.Code)
	}
	return errs.ErrBadRequest.WithMessage(msg)
}
