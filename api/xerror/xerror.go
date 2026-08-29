package xerror

import (
	"fmt"
	"net/http"

	"github.com/morikuni/failure"
)

const (
	CodeClientValidation failure.StringCode = "A-1-10001"
	CodeInvalidChatState failure.StringCode = "E-1-00001"
	CodeUnknownServer    failure.StringCode = "E-0-00001"
)

func ClientValidationErr(err error) error {
	if err == nil {
		return nil
	}
	return failure.Translate(err, CodeClientValidation, failure.Message("invalid request"))
}

func IsCustomError(err error) bool {
	_, ok := failure.CodeOf(err)
	return ok
}

func InvalidChatState(state any) error {
	return failure.New(CodeInvalidChatState, failure.Message(fmt.Sprintf("unknown chat session state: %s", state)))
}

func UnknownServerErr(err error) error {
	if err == nil {
		return nil
	}
	return failure.Translate(err, CodeUnknownServer, failure.Message("internal server error"))
}

func HTTPStatus(err error) int {
	code, ok := failure.CodeOf(err)
	if !ok {
		return http.StatusInternalServerError
	}
	if code == CodeClientValidation {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	message, ok := failure.MessageOf(err)
	if !ok {
		return err.Error()
	}
	return message
}

func Code(err error) string {
	code, ok := failure.CodeOf(err)
	if !ok {
		return "-1"
	}
	return code.ErrorCode()
}
