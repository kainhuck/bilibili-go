package main

import (
	"fmt"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

func main() {
	client := bilibili_go.NewClient()
	client.LoginWithQrCodeWithCache()

	resp, err := client.SubmitVideo("/Users/edy/Downloads/一起去郊游吧.mp4", "/Users/edy/Downloads/111.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("success", resp.Aid, resp.Bvid)
}
