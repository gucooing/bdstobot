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
		logger.Warn("连接外置 discord bot 失败:", err)
		return
	}
	defer func() {
		if err := conndiscordbot.Close(); err != nil {
		}
	}()
	logger.Info("外置 discord bot ws 连接成功")
	for {
		_, message, err := conndiscordbot.ReadMessage()
		if err != nil {
			logger.Warn("监听外置 discord bot 消息失败:", err)
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
	logger.Debug("接收外置 discord bot 消息:", string(message))
	// 解析JSON
	var rsab Rsab
	err := json.Unmarshal(message, &rsab)
	if err != nil {
		logger.Warn("解析JSON失败:", err)
		return ""
	}

	data := decrypt.Protoxor(rsab.Content, rsab.Sign)
	newdata, _ := base64.StdEncoding.DecodeString(data)
	// 使用proto.Unmarshal函数将字节切片反序列化为Person对象
	newPerson := &proto2.Discordbot{}
	err = proto.Unmarshal(newdata, newPerson)
	if err != nil {
		logger.Warn("反序列化失败:", err)
		return "反序列化失败"
	}
	logger.Debug("反protobuf序列化的结果是:", newPerson)

	if newPerson.Type == "cmd" {
		logger.Debug("discord用户名:", newPerson.User)
		takeover.Pflpwsreq(newPerson.Type, newPerson.Cause)
	}
	return ""
}
