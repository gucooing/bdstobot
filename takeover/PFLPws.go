package takeover

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/encryption"
	jsoniter "github.com/json-iterator/go"
)

var connpflp *websocket.Conn

// Playe 定义接收结构体
type Playes struct {
	Type   string   `json:"type"`
	Cause  string   `json:"cause"`
	Action string   `json:"action"`
	Params *Paramss `json:"params"`
}

type Paramss struct {
	Sender string `json:"sender"`
	Xuid   string `json:"xuid"`
	Ip     string `json:"ip"`
	Cmd    string `json:"cmd"`
	Id     string `json:"id"`
	Text   string `json:"text"`
}

// SendWSMessage 定义发送函数
func sendWSMessage(msg []byte) error {
	// 检查是否已经存在连接
	if connpflp == nil {
		serverURL := config.GetConfig().PFLPWsurl
		var err error
		connpflp, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// 发送消息
	err := connpflp.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SendWSMessagesi 定义发送函数
func Pflpwsreq(types, msg string) {
	if types == "cmd" {
		playe := Playes{
			Type:   "pack",
			Action: "runcmdrequest",
			Params: &Paramss{
				Cmd: msg,
				Id:  "0",
			},
		}
		jpkt, _ := jsoniter.Marshal(playe)
		newplaye := encryption.Encrypt_send(string(jpkt))
		// 发送消息
		fmt.Printf("向 PFLP发送 发送数据: %v\n", string(newplaye))
		err := sendWSMessage(newplaye)
		if err != nil {
			return
		}
		return
	}
	if types == "chat" {
		playe := Playes{
			Type:   "pack",
			Action: "sendtext",
			Params: &Paramss{
				Text: msg,
				Id:   "0",
			},
		}
		jpkt, _ := jsoniter.Marshal(playe)
		newplaye := encryption.Encrypt_send(string(jpkt))
		// 发送消息
		fmt.Printf("向 PFLP发送 发送数据: %v\n", string(newplaye))
		err := sendWSMessage(newplaye)
		if err != nil {
			return
		}
	}
}
