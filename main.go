package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal"
	"github.com/gucooing/bdstobot/internal/http"
	"github.com/gucooing/bdstobot/pkg/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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
	cfg := config.GetConfig()
	httpsrv := http.NewServer(cfg)
	if httpsrv == nil {
		fmt.Print("服务器初始化失败")
		return
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpsrv.Start(); err != nil {
			logger.Error("无法启动HTTP服务器")
		}
	}()

	restartTicker := time.NewTicker(time.Duration(999999999) * time.Second)
	go func() {
		for {
			select {
			case <-restartTicker.C:
				logger.Info("正在重启服务器...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := httpsrv.Shutdown(ctx); err != nil {
					logger.Error("无法正常关闭HTTP服务器")
				}
				// 在这里做任何需要的清理工作
				err := http.Restart()
				if err != nil {
					logger.Error("无法重启服务器")
				}
			case <-done:
				// 添加停止服务
				restartTicker.Stop()
				logger.Info("HTTP服务正在关闭")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := httpsrv.Shutdown(ctx); err != nil {
					logger.Error("无法正常关闭HTTP服务")
				}
				logger.Info("HTTP服务已停止")
				os.Exit(0) // 将终止程序

			}
		}
	}()

	internal.Start()
}
