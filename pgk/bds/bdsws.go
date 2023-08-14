package bds

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk"
	"strconv"
	"time"
)

var conn *websocket.Conn = nil

// Reqws 函数用于建立与 cqhttp 的 WebSocket 连接
func Reqws() {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().BdsWsurl
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
		}
	}()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		_ = reswsdata(string(message))
	}
}

// Playe 定义接收结构体
type Playe struct {
	Type   string  `json:"type"`
	Cause  string  `json:"Cause"`
	Params *Params `json:"Params"`
}

type Params struct {
	Sender string `json:"Sender"`
	Xuid   string `json:"Xuid"`
	Ip     string `json:"Ip"`
	Text   string `json:"Text"`
}

func reswsdata(message string) string {
	// 解析JSON
	var playe Playe
	err := json.Unmarshal([]byte(message), &playe)
	if err != nil {
		fmt.Println("解析 JSON 出错:", err)
		return ""
	}
	times := time.Now().Unix()
	fmt.Printf("ws接收数据: %v\n", message)
	if playe.Cause == "join" {
		msg := "玩家：" + playe.Params.Sender + " 偷偷的加入服务器.(<t:" + strconv.Itoa(int(times)) + ":R>)"
		fmt.Printf("发送数据: %v\n", msg)
		pgk.Discord(msg)
	}
	if playe.Cause == "left" {
		msg := "玩家：" + playe.Params.Sender + " 悄悄地退出服务器.(<t:" + strconv.Itoa(int(times)) + ":R>)"
		pgk.Discord(msg)
	}
	if playe.Cause == "chat" {
		msg := "玩家：" + playe.Params.Sender + " 说：" + playe.Params.Text + "(<t:" + strconv.Itoa(int(times)) + ":R>)"
		pgk.Discord(msg)
	}
	return ""
}

// SendWSMessage 定义发送函数
func SendWSMessage(msg interface{}) error {
	var err error
	serverURL := config.GetConfig().CqhttpWsurl
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	// 发送消息
	err = conn.WriteJSON(msg)
	if err != nil {
		return err
	}
	return nil
}
