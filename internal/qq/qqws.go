package qq

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/danger"
	"github.com/gucooing/bdstobot/internal/db"
	"github.com/gucooing/bdstobot/internal/dealwith"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/gucooing/bdstobot/pkg/state"
	"github.com/gucooing/bdstobot/takeover"
	"regexp"
	"strconv"
	"time"
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
	serverURL := config.GetConfig().QQ.CqhttpWsurl
	connqq, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		logger.Warn("连接接收cqhttp失败:", err)
		return
	}
	defer connqq.Close()
	logger.Info("cqhttp ws 连接成功")
	go func() {
		for {
			// 检查是否已经存在连接
			if connqq == nil {
				return
			}
			// 创建并发送 ping 消息
			err := connqq.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				logger.Warn("qq bot ping 发送失败")
				return
			}

			time.Sleep(30 * time.Second)
		}
	}()
	for {
		// 检查是否已经存在连接
		if connqq == nil {
			logger.Warn("cqhttp ws 连接断开，退出循环返回主函数重连")
			return
		}

		_, message, err := connqq.ReadMessage()
		if err != nil {
			logger.Error("cqhttp ws 接收失败,错误：", err)
			return
		}
		go reswsdata(message)
	}
}

func reswsdata(message []byte) {
	// 解析JSON
	logger.Debug("接收 cqhttp ws 消:", string(message))
	var cqhttppost Cqhttppost
	err := json.Unmarshal(message, &cqhttppost)
	if err != nil {
		logger.Warn("解析JSON失败:", err)
		return
	}

	switch cqhttppost.Message {
	case "mc 启动!":
		if cqhttppost.UserId == config.GetConfig().QQ.QqAdmin {
			back := danger.Cmdstart("start " + config.GetConfig().Mcpath)
			takeover.Wscqhttpreq(back)
			return
		} else {
			takeover.Wscqhttpreq("您不是管理员！")
			return
		}
	case "解绑":
		if cqhttppost.GroupId == config.GetConfig().QQ.QQgroup {
			name := db.FindGameNameByQQ(cqhttppost.UserId)
			if name != "" {
				logger.Debug("解绑的游戏昵称为:", name)
				dealwith.Untie(cqhttppost.UserId, name)
				return
			} else {
				takeover.Wscqhttpreq("您没有绑定")
				return
			}
		}
	case "服务器状态":
		motddata, err := state.MotdBE(config.GetConfig().McHost)
		if err != nil {
			logger.Warn("获取motd状态失败 错误:", err)
			return
		}
		msg := "服务器版本:" + motddata.Version + "\n服务器支持的协议:" + strconv.Itoa(motddata.Agreement) + "\n在线玩家:" + strconv.Itoa(motddata.Online) + "\n服务器延迟:" + strconv.FormatInt(motddata.Delay, 10) + "\n内存使用情况:" + state.GetMemPercents() + "%\n内存使用量:" + state.GetMemPercent() + "\ncpu使用情况：" + state.GetCpuPercent() + "%\n"
		logger.Debug("获取服务器状态:", msg)
		takeover.Wscqhttpreq(msg)
		return
	}

	//绑定
	re := regexp.MustCompile(`^绑定\s+(.*)$`)
	matches := re.FindStringSubmatch(cqhttppost.Message)
	if len(matches) > 1 {
		if cqhttppost.GroupId == config.GetConfig().QQ.QQgroup {
			logger.Debug("绑定的游戏昵称为:", matches[1])
			dealwith.Tobind(cqhttppost.UserId, matches[1])
			return
		}
	}

	//聊天转发
	res := regexp.MustCompile(`chat([^/]+)$`)
	match := res.FindStringSubmatch(cqhttppost.Message)
	if len(match) > 1 {
		result := match[1]
		logger.Debug("接收QQ群聊转发消息:", result)
		chat := "[" + cqhttppost.Sender.Nickname + "]QQ群聊消息：" + match[1]
		msg := takeover.Pflpwsreq("chat", chat)
		takeover.Wscqhttpreq(msg)
		return
	}

	//管理员发送指令
	qqadmin := regexp.MustCompile(`cmd\s([^/]+)$`)
	qqadmins := qqadmin.FindStringSubmatch(cqhttppost.Message)
	if len(qqadmins) > 1 {
		if cqhttppost.UserId == config.GetConfig().QQ.QqAdmin {
			//takeover.Wscqhttpreq("执行成功！")
			msg := takeover.Pflpwsreq("cmd", qqadmins[1])
			takeover.Wscqhttpreq(msg)
			return
		} else {
			takeover.Wscqhttpreq("您不是管理员！")
			return
		}
	}

	//管理员发送命令
	qqadminsh := regexp.MustCompile(`sh\s([^/]+)$`)
	qqadminssh := qqadminsh.FindStringSubmatch(cqhttppost.Message)
	if len(qqadminssh) > 1 {
		if cqhttppost.UserId == config.GetConfig().QQ.QqAdmin {
			//takeover.Wscqhttpreq("执行成功！")
			msg := danger.Cmdstart(qqadminssh[1])
			if msg == "" {
				takeover.Wscqhttpreq("错误的，无法执行的命令行")
				return
			}
			takeover.Wscqhttpreq(msg)
			return
		} else {
			takeover.Wscqhttpreq("您不是管理员！")
			return
		}
	}

	return
}
