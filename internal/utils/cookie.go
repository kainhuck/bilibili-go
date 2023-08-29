package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func SaveCookiesToFile(filename string, cookies []*http.Cookie) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, cookie := range cookies {
		_, err := fmt.Fprintf(file, "%s=%s;%s\n", cookie.Name, cookie.Value, cookie.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadCookiesFromFile(filename string) ([]*http.Cookie, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	lines := strings.Split(string(data), "\n")
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
