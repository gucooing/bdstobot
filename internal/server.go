package internal

import (
	"encoding/json"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/bds"
	"github.com/gucooing/bdstobot/internal/discord"
	"github.com/gucooing/bdstobot/internal/discordbot"
	"github.com/gucooing/bdstobot/internal/http"
	"github.com/gucooing/bdstobot/internal/qq"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/gucooing/bdstobot/pkg/state"
	"time"
)

var errorCount int
var nerrorCount int

func Start() {
	go func() {
		http.Httpserver()
		logger.Warn("http服务意外退出，将在5秒后重启")
		time.Sleep(5 * time.Second)
	}()
	//qq部分
	if config.GetConfig().QQ.QQ { //是否启用QQ
		logger.Info("开启使用cqhttp连接QQ")
		go func() {
			for {
				qq.Reqws()
				logger.Warn("cqhttp 失去连接 重连中 ...")
				time.Sleep(5 * time.Second)
			}
		}()
	}
	//discord bot部分
	if config.GetConfig().Discord.DiscordBot {
		logger.Info("使用内置 discord bot")
		go func() {
			for {
				discordbot.DiscordBot()
				logger.Warn("discord bot 失去连接 重连中 ...")
				time.Sleep(5 * time.Second)
			}
		}()
	} else {
		//连接外置discord bot
		logger.Info("使用外置 discord bot")
		go func() {
			for {
				discord.Reqws()
				logger.Warn("discord bot 失去连接 重连中 ...")
				time.Sleep(5 * time.Second)
			}
		}()
	}
	go func() { //连接PFLP ws
		for {
			bds.Reqws()
			logger.Warn("与bds服务器插件 PFLP 失去连接 10秒后将尝试重连 ...")
			time.Sleep(10 * time.Second)
		}
	}()
	for { //死循环保活+服务器状态监控
		data, err := state.MotdBE(config.GetConfig().McHost)
		if errorCount == 3 {
			bds.Nreswsdata("bds服务器掉线 尝试重连")
			logger.Warn("bds服务器掉线 尝试重连")
			nerrorCount = 1
		}
		if err != nil {
			errorCount++
			logger.Warn("获取motd状态失败 错误:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if nerrorCount == 1 {
			bds.Nreswsdata("bds服务器重连成功")
			logger.Info("bds服务器重连成功")
		}
		//datazm, err := state.MotdPm(config.GetConfig().ZmHost)
		//datajsonzm, _ := json.Marshal(datazm)
		//logger.Info("motd回调:", string(datajsonzm))
		datajson, _ := json.Marshal(data)
		logger.Debug("motd回调:", string(datajson))
		errorCount = 0
		nerrorCount = 0
		time.Sleep(5 * time.Second)
	}
}
