package takeover

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/encryption"
	"github.com/gucooing/bdstobot/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"time"
)

// Playes 定义接收结构体
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
	Result string `json:"result"`
}

func Pflpwsreq(types, msg string) string {
	serverURL := config.GetConfig().Pflp.PFLPWsurl
	connpflps, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		logger.Warn("发送pflpws连接失败，错误：", err)
		return ""
	}
	defer connpflps.Close()
	logger.Debug("发送pflp ws 连接成功")
	msgg := pflptype(types, msg)
	if msgg == nil {
		logger.Warn("pflp发送消息处理失败")
		return ""
	}
	// 发送消息
	logger.Debug("向 PFLP发送 发送数据:", string(msgg))
	// 发送消息
	err = connpflps.WriteMessage(websocket.TextMessage, msgg)
	if err != nil {
		logger.Warn("发送PFLP ws 消息失败:", err)
		return ""
	}
	logger.Debug("发送PFLP ws 消息成功")
	time.Sleep(1 * time.Second)
	_, message, err := connpflps.ReadMessage()
	if err != nil {
		logger.Warn("接收PFLP ws 消息失败:", err)
		return ""
	}
	logger.Debug("接收PFLP ws 消息成功:", message)
	var playe Playes
	err = json.Unmarshal([]byte(message), &playe)
	if err != nil {
		logger.Warn("解析 JSON 出错:", err)
		return ""
	}
	return playe.Params.Result
}

// 通过type判断如何处理消息内容
func pflptype(types, msg string) []byte {
	switch types {
	case "cmd":
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
		return newplaye
	case "chat":
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
		return newplaye
	}
	return nil
}
