package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal"
	"github.com/gucooing/bdstobot/pkg/logger"
	"os"
	"strings"
)

func main() {
	// 启动读取配置
	err := config.LoadConfig()
	if err != nil {
		if err == config.FileNotExist {
			p, _ := json.MarshalIndent(config.DefaultConfig, "", "  ")
			fmt.Printf("找不到配置文件，这是默认配置:\n%s\n", p)
			fmt.Printf("\n您可以将其保存到名为“config.json”的文件中并再次运行该程序\n")
			fmt.Printf("按 'Enter' 键退出 ...\n")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		} else {
			panic(err)
		}
	}
	// 初始化日志
	logger.InitLogger()
	logger.SetLogLevel(strings.ToUpper(config.GetConfig().LogLevel))
	logger.Info("    ___     __    _____    __")
	logger.Info("   /   |   / /   / ___/   / /")
	logger.Info("  / /| |  / /    \\__ \\   / /")
	logger.Info(" / ___ | / /___ ___/ /  / /___")
	logger.Info("/_/  |_|/_____//____/  /_____/")

	// 启动服务器
	//test := "{\"Name\":\"1872507219\",\"GameName\":\"xlpmyxhdr\"}"
	//db.Mysqladd(test)
	internal.Start()
}
