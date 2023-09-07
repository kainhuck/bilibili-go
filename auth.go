package bilibili_go

import (
	"encoding/json"
	"net/http"
	"os"
)

type AuthInfo struct {
	Cookies      []*http.Cookie `json:"cookies"`
	RefreshToken string         `json:"refresh_token"`
}

type AuthStorage interface {
	// LoadAuthInfo 加载AuthInfo
	LoadAuthInfo() (*AuthInfo, error)

	// SaveAuthInfo 保存AuthInfo
	SaveAuthInfo(*AuthInfo) error
}

type fileAuthStorage struct {
	file string
}

func (f fileAuthStorage) LoadAuthInfo() (*AuthInfo, error) {
	bts, err := os.ReadFile(f.file)
	if err != nil {
		return nil, err
	}

	var auth AuthInfo
	if err = json.Unmarshal(bts, &auth); err != nil {
		return nil, err
	}

	return &auth, nil
}

func (f fileAuthStorage) SaveAuthInfo(info *AuthInfo) error {
	bts, err := json.Marshal(info)
	if err != nil {
		return err
	}

	file, err := os.Create(f.file)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bts)

	return err
}

func NewFileAuthStorage(file string) AuthStorage {
	return &fileAuthStorage{file: file}
}
