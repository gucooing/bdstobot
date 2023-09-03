package takeover

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/logger"
	"time"
)

// Paramsqq 定义发送结构体
type Paramsqq struct {
	//MessageType string `json:"message_type"`
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type Rsqdataqq struct {
	Action string    `json:"action"`
	Params *Paramsqq `json:"params"`
}

func Wscqhttpreq(msg string) {
	serverURL := config.GetConfig().CqhttpWsurl
	connqq, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		logger.Warn("连接 cqhttp ws 失败:", err)
		return
	}
	defer connqq.Close()
	logger.Debug("发送 cqhttp ws 连接成功")
	msgg := cqhttpmsg(msg)
	if msgg == nil {
		logger.Warn("发送消息处理失败")
		return
	}
	logger.Debug("向 cqhttp ws发送 发送数据:", string(msgg))
	err = connqq.WriteMessage(websocket.TextMessage, msgg)
	if err != nil {
		logger.Warn("发送cqhttp ws 消息失败:", err)
		return
	}
	logger.Debug("发送cqhttp ws 消息成功")
	time.Sleep(3 * time.Second)
	_, message, err := connqq.ReadMessage()
	if err != nil {
		logger.Warn("接收cqhttp ws 消息失败:", err)
		return
	}
	logger.Debug("接收cqhttp ws 消息成功:", message)
	return
}

func cqhttpmsg(msg string) []byte {
	rsqdata := Rsqdataqq{
		Action: "send_group_msg",
		Params: &Paramsqq{
			//MessageType: "private",
			GroupId:    config.GetConfig().QQgroup,
			Message:    msg,
			AutoEscape: false,
		},
	}
	// 发送消息
	reqdatajson, err := json.Marshal(rsqdata)
	if err != nil {
		return nil
	}
	return reqdatajson
}
