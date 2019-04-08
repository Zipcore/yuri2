package api

import (
	"encoding/json"
	"io"
	"net/http"
)

var stdErrMsgs = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	429: "too many requests",
}

type apiErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errResponse(w http.ResponseWriter, code int, msg string) {
	if msg == "" {
		msg = stdErrMsgs[code]
	}

	data := map[string]interface{}{
		"error": apiErrorBody{
			Code:    code,
			Message: msg,
		},
	}

	jsonResponse(w, code, data)
}

func parseJSONBody(body io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(body)
	return dec.Decode(v)
}

func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	var bData []byte
	var err error

	w.Header().Add("Content-Type", "application/json")

	if data != nil {
		bData, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(bData)
}