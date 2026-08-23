package response

import (
	"adventuria/internal/adventuria/errs"
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

type result struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitzero"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func getLang(e *core.RequestEvent) string {
	lang := e.Request.Header.Get("Accept-Language")
	if lang == "" {
		lang = "ru"
	}
	return lang
}

func Error(e *core.RequestEvent, err error) error {
	lang := getLang(e)

	if appErr, ok := errors.AsType[*errs.AppError](err); ok {
		res := result{
			Success: false,
			Data:    nil,
			Error:   appErr.Code,
			Message: appErr.Message,
		}

		status := http.StatusInternalServerError
		if appErr.Status > 0 {
			status = appErr.Status
		}

		if msg, ok := appErr.Translates[lang]; ok {
			res.Message = msg
		}

		return e.JSON(status, res)
	}

	if err := e.JSON(http.StatusInternalServerError, result{
		Success: false,
		Error:   "internal_server_error",
	}); err != nil {
		return err
	}

	return err
}

func Success(e *core.RequestEvent, data any) error {
	return e.JSON(http.StatusOK, result{
		Success: true,
		Data:    data,
	})
}
