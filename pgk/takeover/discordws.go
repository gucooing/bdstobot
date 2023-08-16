package takeover

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/encryption"
	jsoniter "github.com/json-iterator/go"
)

//神b智障方法，临时解决方案

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
		newplaye := encryption.Encrypt_send(string(jpkt))
		// 发送消息
		fmt.Printf("向 PFLP发送 发送数据: %v\n", string(newplaye))
		err := SendWSMessage(newplaye)
		if err != nil {
			return
		}
	}
}
