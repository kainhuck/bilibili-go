package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

func DumpCookies(cookies []*http.Cookie) (string, error) {
	var cookieDump []byte
	buffer := bytes.NewBuffer(cookieDump)

	for _, cookie := range cookies {
		_, err := fmt.Fprintf(buffer, "%s=%s;%s\n", cookie.Name, cookie.Value, cookie.String())
		if err != nil {
			return "", err
		}
	}

	return buffer.String(), nil
}

func LoadCookies(cookieDump string) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	lines := strings.Split(cookieDump, "\n")
	for _, line := range lines {
		cookie := parseCookie(line)
		if cookie != nil {
			cookies = append(cookies, cookie)
		}
	}

	return cookies, nil
}

func parseCookie(cookieStr string) *http.Cookie {
	parts := strings.SplitN(cookieStr, ";", 2)
	if len(parts) >= 1 {
		nameValue := strings.SplitN(parts[0], "=", 2)
		if len(nameValue) == 2 {
			return &http.Cookie{Name: nameValue[0], Value: nameValue[1]}
		}
	}
	return nil
}
