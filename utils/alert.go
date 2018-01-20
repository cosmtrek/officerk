package utils

import (
	"bytes"
)

// SendAlert ...
func SendAlert(data *bytes.Buffer) {
	client := NewHTTPClient()
	// TODO: service url
	client.Post("http://", "application/json", data)
}
