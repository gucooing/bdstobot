package qq

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
)

var conn *websocket.Conn = nil

// Reqws 函数用于建立与 cqhttp 的 WebSocket 连接
func Reqws() {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().CqhttpWsurl
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
	var data map[string]interface{}
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return ""
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

// Params 定义发送结构体
type Params struct {
	//MessageType string `json:"message_type"`
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type Rsqdata struct {
	Action string  `json:"action"`
	Params *Params `json:"params"`
}

// SendWSMessagesi 定义群聊发送函数
func SendWSMessagesi(msg string) {
	//fmt.Println(config.GetConfig().QqAdmin)
	rsqdata := Rsqdata{
		Action: "send_group_msg",
		Params: &Params{
			//MessageType: "private",
			GroupId:    config.GetConfig().QQgroup,
			Message:    msg,
			AutoEscape: false,
		},
	}
	// 发送消息
	fmt.Printf("发送QQ群聊数据: %v\n", rsqdata)
	err := SendWSMessage(rsqdata)
	if err != nil {
		return
	}
}
