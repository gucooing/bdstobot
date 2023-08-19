package internal

import (
	"encoding/json"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/bds"
	"github.com/gucooing/bdstobot/internal/discord"
	"github.com/gucooing/bdstobot/internal/discordbot"
	"github.com/gucooing/bdstobot/internal/qq"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/gucooing/bdstobot/pkg/motd"
	"github.com/gucooing/bdstobot/takeover"
	"time"
)

var errorCount int
var nerrorCount int

func Start() {
	//qq部分
	if config.GetConfig().QQ { //是否启用QQ
		logger.Info().Msg("开启使用cqhttp连接QQ\n")
		go func() {
			for {
				qq.Reqws()
				logger.Warn().Msg("cqhttp 失去连接 重连中 ...\n")
				time.Sleep(5 * time.Second)
			}
		}()
	}
	//discord bot部分
	if config.GetConfig().DiscordBot {
		logger.Info().Msg("使用内置 discord bot\n")
		go func() {
			for {
				discordbot.DiscordBot()
				logger.Warn().Msg("discord bot 失去连接 重连中 ...\n")
				time.Sleep(5 * time.Second)
			}
		}()
	} else {
		//连接外置discord bot
		logger.Info().Msg("使用外置 discord bot\n")
		go func() {
			for {
				discord.Reqws()
				logger.Warn().Msg("discord bot 失去连接 重连中 ...\n")
				time.Sleep(5 * time.Second)
			}
		}()
	}
	go func() { //连接PFLP ws
		for {
			bds.Reqws()
			logger.Warn().Msg("与bds服务器插件 PFLP 失去连接 10秒后将尝试重连 ...\n")
			time.Sleep(10 * time.Second)
		}
	}()
	go func() { //discord rich
		for {
			takeover.Pflpwsres()
		}
	}()
	for { //死循环保活+服务器状态监控
		data, err := motd.MotdBE(config.GetConfig().Host)
		if errorCount == 3 {
			bds.Nreswsdata("bds服务器掉线 尝试重连")
			logger.Warn().Msg("bds服务器掉线 尝试重连")
			nerrorCount = 1
		}
		if err != nil {
			errorCount++
			logger.Warn().Msgf("获取motd状态失败 错误：%d", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if nerrorCount == 1 {
			bds.Nreswsdata("bds服务器重连成功")
			logger.Warn().Msg("bds服务器重连成功")
		}
		datajson, _ := json.Marshal(data)
		logger.Debug().Msg(string(datajson))
		errorCount = 0
		nerrorCount = 0
		time.Sleep(5 * time.Second)
	}
}
