package qq

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/danger"
	"github.com/gucooing/bdstobot/internal/db"
	"github.com/gucooing/bdstobot/internal/dealwith"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/gucooing/bdstobot/takeover"
	"regexp"
)

var connqq *websocket.Conn = nil

type Cqhttppost struct {
	UserId  int64  `json:"user_id"`
	GroupId int64  `json:"group_id"`
	Message string `json:"message"`
	Sender  *sender
}

type sender struct {
	Card     string `json:"card"`
	Nickname string `json:"nickname"`
}

// Reqws 函数用于建立与 cqhttp 的 WebSocket 连接
func Reqws() {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().CqhttpWsurl
	connqq, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		logger.Warn().Msgf("连接接收cqhttp失败：%d", err)
		return
	}
	defer func() {
		if err := connqq.Close(); err != nil {
			logger.Warn().Msgf("连接接收cqhttp失败：%d", err)
			return
		}
	}()
	logger.Info().Msg("cqhttp ws 连接成功")
	for {
		_, message, err := connqq.ReadMessage()
		if err != nil {
			return
		}
		reswsdata(message)
	}
}

func reswsdata(message []byte) {
	// 解析JSON
	logger.Debug().Msgf("接收 cqhttp ws 消息：%d", string(message))
	var cqhttppost Cqhttppost
	err := json.Unmarshal(message, &cqhttppost)
	if err != nil {
		logger.Warn().Msgf("解析JSON失败:%d", err)
		return
	}

	if cqhttppost.Message == "mc 启动!" {
		back := danger.Cmdstart("chcp 936 & start " + config.GetConfig().Mcpath)
		takeover.Wscqhttpreq(back)
		return
	}

	//绑定
	re := regexp.MustCompile(`^绑定\s+(.*)$`)
	matches := re.FindStringSubmatch(cqhttppost.Message)
	if len(matches) > 1 {
		if cqhttppost.GroupId == config.GetConfig().QQgroup {
			logger.Debug().Msgf("绑定的游戏昵称为：%d", matches[1])
			dealwith.Tobind(cqhttppost.UserId, matches[1])
			return
		}
	}

	//解绑
	if cqhttppost.Message == "解绑" {
		if cqhttppost.GroupId == config.GetConfig().QQgroup {
			name := db.FindGameNameByQQ(cqhttppost.UserId)
			if name != "" {
				logger.Debug().Msgf("解绑的游戏昵称为：%d", name)
				dealwith.Untie(cqhttppost.UserId, name)
				return
			} else {
				takeover.Wscqhttpreq("您没有绑定")
				return
			}
		}
	}

	//聊天转发
	res := regexp.MustCompile(`chat([^/]+)$`)
	match := res.FindStringSubmatch(cqhttppost.Message)
	if len(match) > 1 {
		result := match[1]
		logger.Debug().Msgf("接收QQ群聊转发消息：%d", result)
		chat := "[" + cqhttppost.Sender.Nickname + "]QQ群聊消息：" + match[1]
		takeover.Pflpwsreq("chat", chat)
	}

	//管理员发送指令
	qqadmin := regexp.MustCompile(`cmd\s([^/]+)$`)
	qqadmins := qqadmin.FindStringSubmatch(cqhttppost.Message)
	if len(qqadmins) > 1 {
		if cqhttppost.UserId == config.GetConfig().QqAdmin {
			//takeover.Wscqhttpreq("执行成功！")
			takeover.Pflpwsreq("cmd", qqadmins[1])
		} else {
			takeover.Wscqhttpreq("您不是管理员！")
		}
	}

	return
}
