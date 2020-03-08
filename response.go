package gopress

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A Response Object is used to construct the response of the HTTP request.
type Response struct {
	http.ResponseWriter
}

// SetHeader adds a new header to the response headers
func (res Response) SetHeader(name string, content string) {
	res.Header().Set(name, content)
}

// SendStatus writes to the Header the status code and the message to the Body
func (res Response) SendStatus(status int) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("X-Content-Type-Options", "nosniff")
	res.WriteHeader(status)
	fmt.Fprintln(res, http.StatusText(status))
}

// SetStatus sets the status to the response header
func (res Response) SetStatus(status int) Response {
	res.WriteHeader(status)

	return res
}

// Send writes plain text to the response Body
func (res Response) Send(text string) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("X-Content-Type-Options", "nosniff")

	fmt.Fprintln(res, text)
}

// SendJSON converts a interface to JSON format and writes it to the response Body
func (res Response) SendJSON(item interface{}) error {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("X-Content-Type-Options", "nosniff")

	parsed, err := json.Marshal(item)
	if err != nil {
		return err
	}

	res.Write(parsed)
	return nil
}
