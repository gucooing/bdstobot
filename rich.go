package main

import (
	"fmt"
	"time"

	"github.com/hugolgst/rich-go/client"
)

// func Discordrich() {
func main() {
	err := client.Login("1007667861780693154")
	if err != nil {
		panic(err)
	}

	now := time.Now()
	err = client.SetActivity(client.Activity{
		State:      "Genshin lmpact Beta",
		Details:    "Step Into a Vast Magical World of Adventure.",
		LargeImage: "gl",
		LargeText:  "Genshin lmpact Beta",
		SmallImage: "glc",
		SmallText:  "Genshin lmpact Official",
		//Party: &client.Party{
		//	ID:         "-1",
		//	Players:    1,
		//	MaxPlayers: 4,
		//},
		Timestamps: &client.Timestamps{
			Start: &now,
		},
		Buttons: []*client.Button{
			&client.Button{
				Label: "Genshin lmpact Official",
				//Url:   "https://github.com/gucooing/bdstobot",
				Url: "https://genshin.hoyoverse.com",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	for {
		fmt.Printf("心跳保活\n")
		time.Sleep(time.Second * 60)
	}
}
