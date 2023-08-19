package takeover

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/logger"
)

var connqq *websocket.Conn = nil

// SendWSMessage 定义发送函数
func wscqhttpws(msg interface{}) error {
	// 检查是否已经存在连接
	if connqq == nil {
		serverURL := config.GetConfig().CqhttpWsurl
		var err error
		connqq, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
		if err != nil {
			logger.Warn().Msgf("连接 cqhttp ws 失败：%d", err)
			return err
		}
	}
	// 发送消息
	err := connqq.WriteJSON(msg)
	if err != nil {
		logger.Warn().Msgf("发送 cqhttp ws 消息失败：%d", err)
		return err
	}
	return nil
}

// Params 定义发送结构体
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

// SendWSMessagesi 定义群聊发送函数
func Wscqhttpreq(msg string) {
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
	reqdatajson, _ := json.Marshal(rsqdata)
	logger.Debug().Msgf("发送QQ群聊数据: %d\n", string(reqdatajson))
	err := wscqhttpws(rsqdata)
	if err != nil {
		logger.Warn().Msgf("发送 cqhttp ws 消息失败：%d", err)
		return
	}
}
