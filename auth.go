package bilibili_go

import (
	"encoding/json"
	"net/http"
	"os"
)

type authInfo struct {
	Cookies      []*http.Cookie `json:"cookies"`
	RefreshToken string         `json:"refresh_token"`
}

func loadAuthInfoFromFile(filepath string) (*authInfo, error) {
	bts, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var auth authInfo
	if err = json.Unmarshal(bts, &auth); err != nil {
		return nil, err
	}

	return &auth, nil
}

func saveAuthInfoToFile(filepath string, auth *authInfo) error {
	bts, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bts)

	return err
}
