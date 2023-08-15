package pgk

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/motd"
	"strconv"
)

var conn *websocket.Conn = nil

// Rsqdata 定义json结构体
type Rsqdata struct {
	Type  string `json:"Type"`
	Cause string `json:"Cause"`
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
		_ = reswsdata(string(message))
	}
}

func reswsdata(message string) string {
	//fmt.Printf("ws接收数据: %v\n", message)
	// 解析JSON
	var rsqdata Rsqdata
	err := json.Unmarshal([]byte(message), &rsqdata)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return ""
	}
	if rsqdata.Type == "ping" {
		data, _ := motd.MotdBE(config.GetConfig().Host)
		SendWSMessagesil("motd", strconv.Itoa(int(data.Delay)))
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
