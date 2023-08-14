package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk"
	"github.com/gucooing/bdstobot/pgk/bds"
	"github.com/gucooing/bdstobot/pgk/discordbot"
	"github.com/gucooing/bdstobot/pgk/motd"
	"os"
	"time"
)

var errorCount int
var nerrorCount int

func main() {
	fmt.Println("    ___     __    _____    __")
	fmt.Println("   /   |   / /   / ___/   / /")
	fmt.Println("  / /| |  / /    \\__ \\   / /")
	fmt.Println(" / ___ | / /___ ___/ /  / /___")
	fmt.Println("/_/  |_|/_____//____/  /_____/")
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
	go func() {
		for {
			discordbot.DiscordBot()
			fmt.Printf("discord bot 失去连接 重连中 ...\n")
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		for {
			bds.Reqws()
			fmt.Printf("与bds服务器失去连接 5秒后将尝试重连 ...\n")
			time.Sleep(5 * time.Second)
		}
	}()
	for {
		data, err := motd.MotdBE(config.GetConfig().Host)
		if errorCount == 2 {
			pgk.Discord("bds服务器掉线 尝试重连")
			nerrorCount = 1
		}
		if err != nil {
			errorCount++
			fmt.Println("获取motd状态失败 错误：", err)
			continue
		}
		if nerrorCount == 1 {
			pgk.Discord("bds服务器重连成功")
		}
		fmt.Println(data)
		errorCount = 0
		nerrorCount = 0
		time.Sleep(5 * time.Second)
	}
}
