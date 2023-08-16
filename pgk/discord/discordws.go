package discord

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/decrypt"
	proto2 "github.com/gucooing/bdstobot/pgk/proto"
	"github.com/gucooing/bdstobot/pgk/takeover"
)

var (
	conn *websocket.Conn = nil
)

// Rsqdata 定义json结构体
type Rsqdata struct {
	Type  string `json:"Type"`
	Cause string `json:"Cause"`
	User  string `json:"User"`
}

// Reqws 函数用于建立与外置 Discord bot 的 WebSocket 连接
func Reqws() {
	// 检查是否已经存在连接
	if conn != nil {
		return
	}

	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().DiscordWsurl
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
		takeover.SendWSMessagesi(newPerson.Type, newPerson.Cause)
	}
	return ""
}

// SendWSMessage 定义发送函数
func WaiSendWSMessagesi(msg interface{}) error {
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
	err := conn.WriteJSON(msg)
	if err != nil {
		return err
	}
	return nil
}

// SendWSMessagesi 定义群聊发送函数
func SendWSMessagesil(types, msg string) {
	rsqdata := Rsqdata{
		Type:  types,
		Cause: msg,
	}
	// 发送消息
	fmt.Printf("向 Discord bot 发送数据: %v\n", rsqdata)
	err := WaiSendWSMessagesi(rsqdata)
	if err != nil {
		return
	}
}
