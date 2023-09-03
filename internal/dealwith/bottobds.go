package dealwith

import (
	"github.com/gucooing/bdstobot/internal/db"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/gucooing/bdstobot/takeover"
)

func Tobind(UserId int64, name string) {
	namedata := db.FindGameNameByQQ(UserId)
	margss := "whitelist add " + name
	if namedata != "" {
		logger.Debug("用户：%v 已绑定", namedata)
		takeover.Wscqhttpreq("您已绑定")
		return
	} else {
		msg := takeover.Pflpwsreq("cmd", margss)
		if msg != "" {
			db.SaveGameInfo(UserId, name)
			logger.Debug("用户：%v 绑定白名单成功！", namedata)
			takeover.Wscqhttpreq("绑定白名单成功！:" + msg)
			return
		} else {
			logger.Warn("用户：%v 绑定白名单失败，错误出现在意外的地方", namedata)
			takeover.Wscqhttpreq("绑定白名单失败，错误出现在意外的地方，请联系管理员确认")
			return
		}
	}
}

func Untie(UserId int64, name string) {
	margss := "whitelist remove " + name
	msg := takeover.Pflpwsreq("cmd", margss)
	if msg != "" {
		db.DeleteGameInfoByQQ(UserId)
		logger.Debug("已获取用户：%v 白名单成功！", UserId)
		takeover.Wscqhttpreq("解绑成功！:" + msg)

		return
	} else {
		logger.Debug("用户：%v 解绑失败！", UserId)
		takeover.Wscqhttpreq("解绑白名单失败，错误出现在意外的地方，请联系管理员确认")
		return
	}
}

func Bottransit() {

}
