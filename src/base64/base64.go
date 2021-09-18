package base64

import (
	"encoding/base64"
	"net/http"
)

func Encode(b []byte) string {
	var base64EncodingPrefix string
	// Determine the content type of the image file
	mimeType := http.DetectContentType(b)
	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64EncodingPrefix = "data:image/jpeg;base64,"
	case "image/jpg":
		base64EncodingPrefix = "data:image/jpeg;base64,"
	case "image/png":
		base64EncodingPrefix = "data:image/png;base64,"
	}
	return base64EncodingPrefix + base64.StdEncoding.EncodeToString(b)
}