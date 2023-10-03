package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// randomStringSource is the source for generating random strings.
const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321_+"

// defaultMaxUpload is the default max upload size (10 mb)
const defaultMaxUpload = 10485760

type Utils struct {
	MaxJSONSize int
	MaxFileSize int
	AllowedFileTypes []string
	AllowUnknownFields bool
}

type JSONResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty`
}

func (u *Utils) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := defaultMaxUpload
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (u *Utils) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (u *Utils) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()
	return u.WriteJSON(w, statusCode, payload)
}

// RandomString returns a random string of letters of length n, using characters specified in randomStringSource.
func (u *Utils) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}
