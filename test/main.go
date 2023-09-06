package main

import (
	"fmt"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

func main() {
	client := bilibili_go.NewClient(bilibili_go.WithCookieFilePath("bilibili_cookie.txt"))
	client.LoginWithQrCode()

	resp, err := client.GetNavigationStatus()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)

	//SubmitVideo(client)
}

// SubmitVideo 视频投稿
func SubmitVideo(client *bilibili_go.Client) {
	// 1. 上传视频
	video, err := client.UploadVideo("/Users/edy/Downloads/一起去郊游吧.mp4")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("视频上传成功")

	// 2. 上传封面
	cover, err := client.UploadCover("/Users/edy/Downloads/cover.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("封面上传成功")

	// 3. 投稿
	result, err := client.SubmitVideo(&bilibili_go.SubmitRequest{
		Cover:     cover.Url,
		Title:     "一起去郊游吧",
		Copyright: 1,
		TID:       229,
		Tag:       "郊游",
		Desc:      "我们一起去郊游吧",
		Recreate:  -1,
		Videos: []*bilibili_go.Video{
			video,
		},
		NoReprint: 1,
		WebOS:     2,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("投稿成功🏅️AV号: %v, BV号: %v\n", result.Aid, result.Bvid)
}
