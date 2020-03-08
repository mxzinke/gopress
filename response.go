package gopress

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A Response Object is used to construct the response of the HTTP request.
type Response struct {
	writer http.ResponseWriter
	status int
}

// SetHeader adds a new header to the response headers
func (res Response) SetHeader(name string, content string) {
	res.writer.Header().Set(name, content)
}

// SendStatus writes to the Header the status code and the message to the Body
func (res Response) SendStatus(status int) {
	res.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.writer.Header().Set("X-Content-Type-Options", "nosniff")
	res.writer.WriteHeader(status)
	fmt.Fprintln(res.writer, http.StatusText(status))
}

// SetStatus sets the status to the response header
func (res Response) SetStatus(status int) Response {
	res.writer.WriteHeader(status)
	res.status = status

	return res
}

// GetStatus returns the status code which is already set
// If non is set yet, it is automatically status 200
func (res Response) GetStatus() int {
	if res.status == 0 {
		return http.StatusOK
	}

	return res.status
}

// Send writes plain text to the response Body
func (res Response) Send(text string) {
	res.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.writer.Header().Set("X-Content-Type-Options", "nosniff")

	fmt.Fprintln(res.writer, text)
}

// SendJSON converts a interface to JSON format and writes it to the response Body
func (res Response) SendJSON(item interface{}) error {
	res.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.writer.Header().Set("X-Content-Type-Options", "nosniff")

	parsed, err := json.Marshal(item)
	if err != nil {
		return err
	}

	res.writer.Write(parsed)
	return nil
}
