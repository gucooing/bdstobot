package takeover

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/encryption"
	"github.com/gucooing/bdstobot/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
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
	Result string `json:"result"`
}

// Reqws 函数用于建立与 PFLP 的 WebSocket 连接
func Pflpwsres() {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().PFLPWsurl
	connpflp, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		logger.Warn("连接 PFLP ws 失败:", err)
		return
	}
	defer func() {
		if err := connpflp.Close(); err != nil {
		}
	}()
	logger.Info("PFLP ws 连接成功")
	go func() {
		for {
			// 检查是否已经存在连接
			if connpflp == nil {
				return
			}
			// 创建并发送 ping 消息
			err := connpflp.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				logger.Warn("pflp bot ping 发送失败")
				return
			}

			time.Sleep(30 * time.Second)
		}
	}()
	for {
		// 检查是否已经存在连接
		if connpflp == nil {
			return
		}

		_, message, err := connpflp.ReadMessage()
		if err != nil {
			logger.Warn("接收PFLP ws 消息失败:", err)
			return
		}
		_ = reswsdata(string(message))
	}
}

func reswsdata(message string) string {
	// 解析JSON
	logger.Debug("接收PFLP ws 消息:", message)
	var playe Playes
	err := json.Unmarshal([]byte(message), &playe)
	if err != nil {
		logger.Warn("解析 JSON 出错:", err)
		return ""
	}
	times := time.Now().Unix()
	//未加密数据解析处理
	if playe.Cause == "runcmdfeedback" {
		cmdresult := playe.Params.Result
		msg := "回调：" + cmdresult + "(<t:" + strconv.Itoa(int(times)) + ":R>)"
		Wscqhttpreq(msg)
	}
	return ""
}

// SendWSMessage 定义发送函数
func sendWSMessage(msg []byte) error {
	// 检查是否已经存在连接
	if connpflp == nil {
		serverURL := config.GetConfig().PFLPWsurl
		var err error
		connpflp, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
		if err != nil {
			logger.Warn("连接 PFLP ws 失败:", err)
			return err
		}
	}
	// 发送消息
	err := connpflp.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		logger.Warn("发送PFLP ws 消息失败:", err)
		return err
	}
	return nil
}

// SendWSMessagesi 定义发送函数
func Pflpwsreq(types, msg string) bool {
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
		logger.Debug("向 PFLP发送 发送数据:", string(newplaye))
		err := sendWSMessage(newplaye)
		if err != nil {
			return false
		}
		return true
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
		logger.Debug("向 PFLP发送 发送数据:", string(newplaye))
		err := sendWSMessage(newplaye)
		if err != nil {
			logger.Warn("发送PFLP ws 消息失败:", err)
			return false
		}
	}
	return false
}
