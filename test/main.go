package main

import (
	"fmt"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

func main() {
	client := bilibili_go.NewClient(
		bilibili_go.WithCookieFilePath("bilibili_cookie.hyk.txt"),
		bilibili_go.WithDebug(false),
	)
	client.LoginWithQrCode()

	SearchUserInfo(client)

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

// SearchUserInfo 查询用户信息
func SearchUserInfo(client *bilibili_go.Client) {
	// 1. 根据mid查询其他用户信息
	card, err := client.GetUserCard("2", true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("用户名：%v，粉丝数：%v\n", card.Card.Name, card.Card.Fans)

	// 2. 查询自身信息
	resp, err := client.GetMyInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("用户名：%v，粉丝数：%v，硬币数：%v\n", resp.Name, resp.Follower, resp.Coins)
}
