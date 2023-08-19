package discord

import (
	"encoding/base64"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/decrypt"
	"github.com/gucooing/bdstobot/pkg/logger"
	proto2 "github.com/gucooing/bdstobot/proto"
	"github.com/gucooing/bdstobot/takeover"
)

var (
	conndiscordbot *websocket.Conn = nil
)

// Reqws 函数用于建立与外置 Discord bot 的 WebSocket 连接
func Reqws() {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().DiscordWsurl
	conndiscordbot, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		logger.Warn().Msgf("连接外置 discord bot 失败：%d", err)
		return
	}
	defer func() {
		if err := conndiscordbot.Close(); err != nil {
		}
	}()
	logger.Info().Msg("外置 discord bot ws 连接成功")
	for {
		_, message, err := conndiscordbot.ReadMessage()
		if err != nil {
			logger.Warn().Msgf("监听外置 discord bot 消息失败：%d", err)
			return
		}
		_ = biswsdata(message)
	}
}

// 定义加密数据结构体
type Rsab struct {
	Content string `json:"content"`
	Sign    string `json:"sign"`
}

func biswsdata(message []byte) string {
	logger.Debug().Msgf("接收外置 discord bot 消息: %d\n", string(message))
	// 解析JSON
	var rsab Rsab
	err := json.Unmarshal(message, &rsab)
	if err != nil {
		logger.Warn().Msgf("解析JSON失败:%d", err)
		return ""
	}

	data := decrypt.Protoxor(rsab.Content, rsab.Sign)
	newdata, _ := base64.StdEncoding.DecodeString(data)
	// 使用proto.Unmarshal函数将字节切片反序列化为Person对象
	newPerson := &proto2.Discordbot{}
	err = proto.Unmarshal(newdata, newPerson)
	if err != nil {
		logger.Warn().Msgf("反序列化失败:%d", err)
		return "反序列化失败"
	}

	if newPerson.Type == "cmd" {
		logger.Debug().Msgf("discord用户名：%d", newPerson.User)
		takeover.Pflpwsreq(newPerson.Type, newPerson.Cause)
	}
	return ""
}
