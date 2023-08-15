package bds

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk"
	"github.com/gucooing/bdstobot/pgk/qq"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

var conn *websocket.Conn

// Reqws 函数用于建立与 cqhttp 的 WebSocket 连接
func Reqws() {
	// 检查是否已经存在连接
	if conn != nil {
		return
	}

	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().PFLPWsurl
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
	Cause  string  `json:"cause"`
	Action string  `json:"action"`
	Params *Params `json:"params"`
}

type Params struct {
	Sender string `json:"sender"`
	Xuid   string `json:"xuid"`
	Ip     string `json:"ip"`
	Cmd    string `json:"cmd"`
	Id     string `json:"id"`
	Text   string `json:"text"`
}

type Encrypt struct {
	Type   string         `json:"type"`
	Params *EncryptParams `json:"params"`
}

type EncryptParams struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

func reswsdata(message string) string {
	// 解析JSON
	//fmt.Printf("ws接收数据: %v\n", message)
	var playe Playe
	err := json.Unmarshal([]byte(message), &playe)
	if err != nil {
		fmt.Println("解析 JSON 出错:", err)
		return ""
	}
	times := time.Now().Unix()
	//未加密数据解析处理
	if playe.Cause == "join" {
		msg := "玩家：" + playe.Params.Sender + " 偷偷的加入服务器.(<t:" + strconv.Itoa(int(times)) + ":R>)"
		fmt.Printf("发送数据: %v\n", msg)
		Nreswsdata(msg)
	}
	if playe.Cause == "left" {
		msg := "玩家：" + playe.Params.Sender + " 悄悄地退出服务器.(<t:" + strconv.Itoa(int(times)) + ":R>)"
		Nreswsdata(msg)
	}
	if playe.Cause == "chat" {
		msg := "玩家：" + playe.Params.Sender + " 说：" + playe.Params.Text + "(<t:" + strconv.Itoa(int(times)) + ":R>)"
		Nreswsdata(msg)
	}
	return ""
}

// 传递逻辑再处理
func Nreswsdata(msg string) {
	if config.GetConfig().QQ {
		//发送QQ消息
		qq.SendWSMessagesi(msg)
	}
	if config.GetConfig().DiscordBot {
		//使用内置discord bot发送消息
		pgk.Discord(msg)
	} else {
		//使用外置discord bot发送消息
		pgk.SendWSMessagesil("chat", msg)
	}
}

// SendWSMessage 定义发送函数
func SendWSMessage(msg []byte) error {
	// 检查是否已经存在连接
	if conn == nil {
		serverURL := config.GetConfig().CqhttpWsurl
		var err error
		conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
		if err != nil {
			return err
		}
	}
	// 发送消息
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

// SendWSMessagesi 定义发送函数
func SendWSMessagesi(types, msg string) {
	if types == "cmd" {
		//newcmd, _ := pgk.Encrypt([]byte(msg))
		playe := Playe{
			Type:   "pack",
			Action: "runcmdrequest",
			Params: &Params{
				Cmd: msg,
				Id:  "0",
			},
		}
		jpkt, _ := jsoniter.Marshal(playe)
		newplaye := pgk.Encrypt_send(string(jpkt))
		// 发送消息
		fmt.Printf("向 PFLP发送 发送数据: %v\n", string(newplaye))
		err := SendWSMessage(newplaye)
		if err != nil {
			return
		}
	}
}
