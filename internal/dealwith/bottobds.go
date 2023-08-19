package dealwith

import (
	"github.com/gucooing/bdstobot/internal/db"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/gucooing/bdstobot/takeover"
)

func Tobind(UserId int64, name string) {
	namedata := db.FindGameNameByQQ(UserId)
	if namedata != "" {
		logger.Debug().Msgf("用户：%d 已绑定", namedata)
		takeover.Wscqhttpreq("您已绑定")
		return
	} else {
		if db.SaveGameInfo(UserId, name) {
			logger.Debug().Msgf("用户：%d 绑定白名单成功！", namedata)
			takeover.Wscqhttpreq("绑定白名单成功！")
			margss := "whitelist add " + name
			takeover.Pflpwsreq("cmd", margss)
			return
		} else {
			logger.Warn().Msgf("用户：%d 绑定白名单失败，错误出现在意外的地方", namedata)
			takeover.Wscqhttpreq("绑定白名单失败，错误出现在意外的地方，请联系管理员确认")
			return
		}
	}
}

func Untie(UserId int64, name string) {
	if db.DeleteGameInfoByQQ(UserId) {
		logger.Debug().Msgf("用户：%d 绑定白名单成功！", UserId)
		takeover.Wscqhttpreq("解绑成功！")
		margss := "whitelist remove " + name
		takeover.Pflpwsreq("cmd", margss)
		return
	}
}

func Bottransit() {

}
