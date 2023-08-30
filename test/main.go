package main

import (
	"github.com/davecgh/go-spew/spew"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

func main() {
	client := bilibili_go.NewClient()
	client.LoginWithQrCodeWithCache()

	preloadResp, err := client.PreUpload("hello.mp4", 10000)
	if err != nil {
		log.Fatal(err)
	}
	//spew.Dump(preloadResp)

	resp, err := client.GetUploadID(preloadResp, 10000)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
}
