package main

import (
	"github.com/davecgh/go-spew/spew"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

func main() {
	client := bilibili_go.NewClient()
	client.LoginWithQrCodeWithCache()

	nav, err := client.PreUpload("asdsad.txt", 100000)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(nav)
}
