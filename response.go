package gopress

// A Response Interface is used to construct the response of the HTTP request.
type Response interface {
	SendStatus(int)
	SendJSON(interface{})
}
