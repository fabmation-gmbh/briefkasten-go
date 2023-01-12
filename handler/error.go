package handler

import (
	"errors"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/fabmation-gmbh/briefkasten-go/handler/rerr"
	"github.com/fabmation-gmbh/briefkasten-go/internal/log"
)

type httpError struct {
	Statuscode int    `json:"status_code"`
	ErrorCode  int    `json:"error_code"`
	Error      string `json:"error"`
}

// ErrorHandler is used to catch errors thrown inside the routes by ctx.Next(err)
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Statuscode defaults to 500
	code := fiber.StatusInternalServerError
	reqID, _ := c.Context().UserValue(fiber.HeaderXRequestID).(string)

	errCode := -1

	l := log.With(zap.String("request_id", reqID))

	// Check if it's an fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code

		l.Error("Error returned from handler", zap.Error(e))
	} else if errors.Is(err, rerr.Error{}) {
		e := err.(rerr.Error)

		code = e.HTTPCode()
		errCode = int(e.Code)
		logMsg := e.LogMsg()

		if logMsg != "" {
			we := e.Unwrap()

			l.Error(logMsg, zap.Error(we))
		} else if we := e.Unwrap(); we != nil {
			l.Error("Error returned from handler", zap.Error(we))
		} else {
			l.Error("An error occurred while handling the request", zap.Any("rerr_error", e))
		}
	}

	return c.Status(code).JSON(&httpError{
		ErrorCode:  errCode,
		Statuscode: code,
		Error:      err.Error(),
	})
}

// byteSlice2String converts a byte slice to a string in a performant way.
func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
