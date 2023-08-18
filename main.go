package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/bds"
	"github.com/gucooing/bdstobot/internal/discord"
	"github.com/gucooing/bdstobot/internal/discordbot"
	"github.com/gucooing/bdstobot/internal/qq"
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
	//qq部分
	if config.GetConfig().QQ { //是否启用QQ
		fmt.Printf("开启使用cqhttp连接QQ\n")
		go func() {
			for {
				qq.Reqws()
				fmt.Printf("cqhttp 失去连接 重连中 ...\n")
				time.Sleep(5 * time.Second)
			}
		}()
	}
	//discord bot部分
	if config.GetConfig().DiscordBot {
		fmt.Printf("使用内置 discord bot\n")
		go func() {
			for {
				discordbot.DiscordBot()
				fmt.Printf("discord bot 失去连接 重连中 ...\n")
				time.Sleep(5 * time.Second)
			}
		}()
	} else {
		//连接外置discord bot
		fmt.Printf("使用外置 discord bot\n")
		go func() {
			for {
				discord.Reqws()
				fmt.Printf("discord bot 失去连接 重连中 ...\n")
				time.Sleep(5 * time.Second)
			}
		}()
	}
	go func() { //连接PFLP ws
		for {
			bds.Reqws()
			fmt.Printf("与bds服务器插件 PFLP 失去连接 10秒后将尝试重连 ...\n")
			time.Sleep(10 * time.Second)
		}
	}()
	go func() { //discord rich
		for {
			discord.Discordrich()
		}
	}()
	for { //死循环保活+服务器状态监控
		_, err := motd.MotdBE(config.GetConfig().Host)
		if errorCount == 3 {
			bds.Nreswsdata("bds服务器掉线 尝试重连")
			nerrorCount = 1
		}
		if err != nil {
			errorCount++
			fmt.Println("获取motd状态失败 错误：", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if nerrorCount == 1 {
			bds.Nreswsdata("bds服务器重连成功")
		}
		//fmt.Println(data)
		errorCount = 0
		nerrorCount = 0
		time.Sleep(5 * time.Second)
	}
}
