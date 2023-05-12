package utils

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/sha3"
)

func GetRequestFingerprint(r *http.Request) string {
	headers := ""
	headers += r.Header.Get("X-Forwarded-For") + " : "
	headers += r.Header.Get("User-Agent") + ":"
	headers += r.Header.Get("Accept-Language")
	return fmt.Sprintf("%x", sha3.Sum256([]byte(headers)))
}
