package main

import (
	"github.com/davecgh/go-spew/spew"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

func main() {
	client := bilibili_go.NewClient()
	//client.LoginWithQrCode()

	// 获取个人信息
	account, err := client.GetAccount()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(account)

	// 获取导航栏信息
	nav, err := client.GetNavigation()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(nav)
}
