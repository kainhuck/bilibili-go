package main

import bilibili_go "github.com/kainhuck/bilibili-go"

func main() {
	client := bilibili_go.NewClient()
	client.LoginWithQrCode()
}
