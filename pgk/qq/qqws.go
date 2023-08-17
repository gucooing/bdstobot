package qq

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/db"
	"github.com/gucooing/bdstobot/pgk/takeover"
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
		return
	}
	defer func() {
		if err := connqq.Close(); err != nil {
			return
		}
	}()
	fmt.Println("cqhttp ws 连接成功")
	for {
		_, message, err := connqq.ReadMessage()
		if err != nil {
			return
		}
		_ = reswsdata(message)
	}
}

func reswsdata(message []byte) string {
	// 解析JSON
	var cqhttppost Cqhttppost
	err := json.Unmarshal(message, &cqhttppost)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return ""
	}
	re := regexp.MustCompile(`^绑定\s+(.*)$`)
	matches := re.FindStringSubmatch(cqhttppost.Message)
	if len(matches) > 1 {
		if cqhttppost.GroupId == config.GetConfig().QQgroup {
			fmt.Println("绑定的游戏昵称为：", matches[1])
			namedata := db.FindGameNameByQQ(cqhttppost.UserId)
			if namedata != "" {
				takeover.Wscqhttpreq("您已绑定")
				return "用户已绑定"
			} else {
				if db.SaveGameInfo(cqhttppost.UserId, matches[1]) {
					takeover.Wscqhttpreq("绑定白名单成功！")
					margss := "whitelist add " + matches[1]
					takeover.Pflpwsreq("cmd", margss)
					return ""
				} else {
					takeover.Wscqhttpreq("绑定白名单失败，错误出现在意外的地方，请联系管理员确认")
					return ""
				}
			}
		}
	}
	if cqhttppost.Message == "解绑" {
		if cqhttppost.GroupId == config.GetConfig().QQgroup {
			namedata := db.FindGameNameByQQ(cqhttppost.UserId)
			if namedata != "" {
				fmt.Println("解绑的游戏昵称为：", namedata)
				if db.DeleteGameInfoByQQ(cqhttppost.UserId) {
					takeover.Wscqhttpreq("解绑成功！")
					margss := "whitelist remove " + namedata
					takeover.Pflpwsreq("cmd", margss)
					return ""
				}
			} else {
				takeover.Wscqhttpreq("您没有绑定")
				return ""
			}
		}
	}
	res := regexp.MustCompile(`chat([^/]+)$`)
	match := res.FindStringSubmatch(cqhttppost.Message)
	if len(match) > 1 {
		result := match[1]
		fmt.Println(result)
		chat := "[" + cqhttppost.Sender.Nickname + "]QQ群聊消息：" + match[1]
		takeover.Pflpwsreq("chat", chat)
	}
	return "123"
}
