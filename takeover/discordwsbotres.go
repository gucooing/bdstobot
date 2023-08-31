package takeover

import (
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/encryption"
	"github.com/gucooing/bdstobot/pkg/logger"
	jsoniter "github.com/json-iterator/go"
)

var conndiscordbot *websocket.Conn

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
func sendWSMessages(msg []byte) error {
	// 检查是否已经存在连接
	if conndiscordbot == nil {
		serverURL := config.GetConfig().DiscordWsurl
		var err error
		conndiscordbot, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
		if err != nil {
			logger.Warn("连接外置 discord bot 失败:", err)
			return err
		}
	}
	// 发送消息
	err := conndiscordbot.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		logger.Warn("发送外置 discord bot 消息失败:", err)
		return err
	}
	return nil
}

// SendWSMessagesi 定义发送函数
func Discordbotwsreq(types, msg string) {
	if types == "cmd" {
		//newcmd, _ := pkg.Encrypt([]byte(msg))
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
		logger.Debug("向 PFLP发送 发送数据:", string(newplaye))
		err := sendWSMessages(newplaye)
		if err != nil {
			return
		}
	}
}
