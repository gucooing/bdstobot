package main

import (
	"bufio"
	"encoding/json"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal"
	"github.com/gucooing/bdstobot/pkg/logger"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	// 启动读取配置
	err := config.LoadConfig()
	if err != nil {
		if err == config.FileNotExist {
			p, _ := json.MarshalIndent(config.DefaultConfig, "", "  ")
			logger.Warn().Msgf("找不到配置文件，这是默认配置:\n%s\n", p)
			logger.Warn().Msg("\n您可以将其保存到名为“config.json”的文件中并再次运行该程序\n")
			logger.Warn().Msg("按 'Enter' 键退出 ...\n")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		} else {
			panic(err)
		}
	}
	logger.Warn().Msg("    ___     __    _____    __")
	logger.Warn().Msg("   /   |   / /   / ___/   / /")
	logger.Warn().Msg("  / /| |  / /    \\__ \\   / /")
	logger.Warn().Msg(" / ___ | / /___ ___/ /  / /___")
	logger.Warn().Msg("/_/  |_|/_____//____/  /_____/")
	if err != nil {
		panic(err)
	}
	switch config.GetConfig().LogLevel {
	case "trace":
		logger.Logger = logger.Logger.Level(zerolog.TraceLevel)
	case "debug":
		logger.Logger = logger.Logger.Level(zerolog.DebugLevel)
	case "info":
		logger.Logger = logger.Logger.Level(zerolog.InfoLevel)
	case "silent", "disabled":
		logger.Logger = logger.Logger.Level(zerolog.Disabled)
	}
	internal.Start()
}
