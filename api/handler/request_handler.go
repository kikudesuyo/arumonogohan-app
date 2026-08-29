package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/xerror"
)

type ProcessFunc func(r *http.Request, requestData map[string]interface{}) (http.Handler, error)

func HandleRequestAndResponse(r *http.Request, w http.ResponseWriter, processFn ProcessFunc) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	requestData, err := getParameters(r)
	if err != nil {
		entity.NewErrorResponse(err).ServeHTTP(w, r)
		return
	}
	resp, err := processFn(r, requestData)
	if err != nil {
		if !xerror.IsCustomError(err) {
			err = xerror.UnknownServerErr(err)
		}
		entity.NewErrorResponse(err).ServeHTTP(w, r)
		return
	}
	resp.ServeHTTP(w, r)
}

func getParameters(r *http.Request) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for key, values := range r.URL.Query() {
		if len(values) == 1 {
			data[key] = values[0]
		} else {
			data[key] = values
		}
	}

	contentType := r.Header.Get("Content-Type")
	if r.Method != http.MethodPost && r.Method != http.MethodPut && !strings.HasPrefix(contentType, "application/json") {
		return data, nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if len(body) == 0 {
		return data, nil
	}
	var bodyData map[string]interface{}
	if err := json.Unmarshal(body, &bodyData); err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	for key, value := range bodyData {
		data[key] = value
	}
	return data, nil
}
