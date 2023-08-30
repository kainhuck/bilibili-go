package main

import (
	"fmt"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

func main() {
	client := bilibili_go.NewClient()
	client.LoginWithQrCodeWithCache()

	err := client.Upload("/Users/edy/Downloads/whale-宣传视频.mp4")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("success")
}
