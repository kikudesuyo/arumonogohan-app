package entity

import (
	"encoding/json"
	"net/http"

	"github.com/kikudesuyo/arumonogohan-app/api/xerror"
)

type BaseResponse struct {
	Meta     ResponseMeta   `json:"meta"`
	Err      *ResponseError `json:"error,omitempty"`
	DataType string         `json:"data_type"`
	Data     any            `json:"data,omitempty"`
}

type ResponseMeta struct {
	StatusCode int  `json:"status_code"`
	Success    bool `json:"success"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewObjectResponse(data any) BaseResponse {
	return BaseResponse{Meta: ResponseMeta{StatusCode: http.StatusOK, Success: true}, DataType: "object", Data: data}
}

func NewErrorResponse(err error) BaseResponse {
	return BaseResponse{
		Meta:     ResponseMeta{StatusCode: xerror.HTTPStatus(err), Success: false},
		Err:      &ResponseError{Code: xerror.Code(err), Message: xerror.Message(err)},
		DataType: "error",
	}
}

func (r BaseResponse) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Meta.StatusCode)
	_ = json.NewEncoder(w).Encode(r)
}
