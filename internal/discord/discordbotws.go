package discord

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/decrypt"
	proto2 "github.com/gucooing/bdstobot/proto"
	"github.com/gucooing/bdstobot/takeover"
)

var (
	conndiscordbot *websocket.Conn = nil
)

// Reqws 函数用于建立与外置 Discord bot 的 WebSocket 连接
func Reqws() {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().DiscordWsurl
	conndiscordbot, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return
	}
	defer func() {
		if err := conndiscordbot.Close(); err != nil {
		}
	}()
	fmt.Println("外置 discord bot ws 连接成功")
	for {
		_, message, err := conndiscordbot.ReadMessage()
		if err != nil {
			return
		}
		_ = biswsdata(message)
	}
}

// 定义加密数据结构体
type Rsab struct {
	Content string `json:"content"`
	Sign    string `json:"sign"`
}

func biswsdata(message []byte) string {
	//fmt.Printf("ws接收数据: %v\n", message)
	// 解析JSON
	var rsab Rsab
	fmt.Println(message)
	err := json.Unmarshal(message, &rsab)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return ""
	}

	data := decrypt.Protoxor(rsab.Content, rsab.Sign)
	newdata, _ := base64.StdEncoding.DecodeString(data)
	// 使用proto.Unmarshal函数将字节切片反序列化为Person对象
	newPerson := &proto2.Discordbot{}
	err = proto.Unmarshal(newdata, newPerson)
	if err != nil {
		fmt.Println("反序列化失败:", err)
		return "反序列化失败"
	}

	if newPerson.Type == "cmd" {
		fmt.Println("discord用户名：", newPerson.User)
		takeover.Pflpwsreq(newPerson.Type, newPerson.Cause)
	}
	return ""
}
