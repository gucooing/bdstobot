package bds

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/discord"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/gucooing/bdstobot/takeover"
	"strconv"
	"time"
)

var connpflp *websocket.Conn

// Reqws 函数用于建立与 cqhttp 的 WebSocket 连接
func Reqws() {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().PFLPWsurl
	connpflp, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return
	}
	defer func() {
		if err := connpflp.Close(); err != nil {
		}
	}()
	logger.Info().Msg("PFLP ws 连接成功")
	for {
		_, message, err := connpflp.ReadMessage()
		if err != nil {
			return
		}
		_ = reswsdata(string(message))
	}
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

func reswsdata(message string) string {
	// 解析JSON
	logger.Debug().Msgf("接收 PFLP ws 消息: %d\n", message)
	var playe Playe
	err := json.Unmarshal([]byte(message), &playe)
	if err != nil {
		logger.Warn().Msgf("解析 JSON 出错:%d", err)
		return ""
	}
	times := time.Now().Unix()
	//未加密数据解析处理
	if playe.Cause == "join" {
		msg := "玩家：" + playe.Params.Sender + " 偷偷的加入服务器.(<t:" + strconv.Itoa(int(times)) + ":R>)"
		Nreswsdata(msg)
	}
	if playe.Cause == "left" {
		msg := "玩家：" + playe.Params.Sender + " 悄悄地退出服务器.(<t:" + strconv.Itoa(int(times)) + ":R>)"
		Nreswsdata(msg)
	}
	if playe.Cause == "chat" {
		msg := "玩家：" + playe.Params.Sender + " 说：" + playe.Params.Text + "(<t:" + strconv.Itoa(int(times)) + ":R>)"
		Nreswsdata(msg)
	}
	return ""
}

// 传递逻辑再处理
func Nreswsdata(msg string) {
	if config.GetConfig().QQ {
		//发送QQ消息
		takeover.Wscqhttpreq(msg)
		logger.Debug().Msgf("发送QQ消息:%d", msg)
	}
	//发送discord webhook消息
	discord.Discordwebhook(msg)
	logger.Debug().Msgf("发送discord webhook消息:%d", msg)
	//discord bot 主动发送消息 暂时无效
	if config.GetConfig().DiscordBot {
		//使用内置discord bot发送消息
		logger.Debug().Msgf("使用内置discord bot发送消息:%d", msg)
	} else {
		//使用外置discord bot发送消息
		logger.Debug().Msgf("使用外置discord bot发送消息:%d", msg)
		//takeover.Discordbotwsreq("chat", msg)
	}
}
