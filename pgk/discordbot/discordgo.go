package discordbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/motd"
	"strconv"
	"strings"
	"time"
)

const prefix = "!" // 机器人命令前缀

func DiscordBot() {
	dg, err := discordgo.New(config.GetConfig().DiscordBotToken)
	if err != nil {
		fmt.Println("创建discord连接失败: ", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("discord bot 连接失败: ", err)
		return
	}

	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			return
		}
	}(dg)

	for {
		// 阻塞携程，保持机器人在线
		time.Sleep(10 * time.Second)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return // 忽略机器人自身的消息
	}

	if !hasPrefix(m.Content, prefix) {
		return // 忽略非命令消息
	}

	args := splitArgs(m.Content, prefix)
	command := args[0]
	fmt.Printf("监听的消息内容是：%s\n", command)
	switch command {
	case "!chat":
		data, err := motd.MotdBE(config.GetConfig().Host)
		if err != nil {
			fmt.Println("获取motd状态失败 错误：", err)
		}
		SendMessage(
			s,
			&SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "服务器延迟为：" + strconv.Itoa(int(data.Delay)),
			},
		)
	case "!cmd":
		fmt.Println("!cmd")
		return
	}
}

func hasPrefix(s string, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func splitArgs(s string, prefix string) []string {
	return removeEmptyStrings(strings.Split(s, " "))
}

func removeEmptyStrings(strings []string) []string {
	result := []string{}
	for _, str := range strings {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
