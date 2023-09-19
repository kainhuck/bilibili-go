package main

import (
	"encoding/json"
	"fmt"
	bilibili_go "github.com/kainhuck/bilibili-go"
	"github.com/skip2/go-qrcode"
	"log"
)

func main() {
	client := bilibili_go.NewClient(
		bilibili_go.WithAuthStorage(bilibili_go.NewFileAuthStorage("bilibili.json")),
		bilibili_go.WithDebug(false),
		bilibili_go.WithShowQRCodeFunc(func(code *qrcode.QRCode) error {

			return code.WriteFile(640, "qrcode.png")
		}),
	)
	client.LoginWithQrCode()

	if err := client.RefreshAuthInfo(); err != nil {
		log.Fatal(err)
	}

	err := client.UnLikeVideo("BV19u4y1D7GT")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OK")

	printIt(client.GetExpReword())

	//printIt(client.GetFriends())
	//tags, err := client.GetRelationTags()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//for _, tag := range tags {
	//	if tag.Count == 0 {
	//		continue
	//	}
	//	pn := 1
	//	for pn < tag.Count {
	//		printIt(client.GetRelationTagUsers(tag.TagId, "", 1, pn))
	//		pn++
	//	}
	//}

	//RelationDemo(client)
	//SearchUserInfo(client)
	//SubmitVideo(client)
}

// SubmitVideo 视频投稿
func SubmitVideo(client *bilibili_go.Client) {
	// 1. 上传视频
	video, err := client.UploadVideoFromDisk("/Users/edy/Downloads/一起去郊游吧.mp4")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("视频上传成功")

	// 2. 上传封面
	cover, err := client.UploadCoverFromDisk("/Users/edy/Downloads/cover.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("封面上传成功")

	// 3. 投稿
	result, err := client.SubmitVideo(&bilibili_go.SubmitRequest{
		Cover:     cover.Url,
		Title:     "一起去郊游吧",
		Copyright: 1,
		TID:       bilibili_go.LifeGroup.RandomTid(),
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
	card, err := client.GetUserCard("13868000", true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("用户名：%v，粉丝数：%v，头衔：%v\n", card.Card.Name, card.Card.Fans, card.Card.Official.Title)

	// 2. 查询自身信息
	resp, err := client.GetMyInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("用户名：%v，粉丝数：%v，硬币数：%v\n", resp.Name, resp.Follower, resp.Coins)
}

// RelationDemo 关系操作
func RelationDemo(client *bilibili_go.Client) {
	// 1. 查询自己的所有粉丝
	pn := 0
	for {
		pn++
		resp, err := client.GetFollowers(50, pn)
		if err != nil {
			log.Fatal(err)
		}

		if len(resp.List) == 0 {
			break
		}

		for _, each := range resp.List {
			// 查询粉丝详细信息
			user, err := client.GetUserInfo(each.Mid)
			if err != nil {
				log.Println(each.Uname, err)
				continue
			}
			fmt.Printf("名字: %v\tmid: %v\t性别: %v\t学校: %v\t等级: %v\n", user.Name, user.Mid, user.Sex, user.School.Name, user.Level)
		}
	}
}

func printIt(a any, err error) {
	if err != nil {
		log.Fatal(err)
	}
	bts, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println(string(bts))
}
