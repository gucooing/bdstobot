package discordbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/motd"
	"strconv"
	"time"
)

var BotId string

func DiscordBot() {
	goBot, err := discordgo.New("Bot " + config.GetConfig().DiscordBotToken)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("discord bot 已成功连接 !")
	for {
		// 阻塞携程，保持机器人在线
		time.Sleep(10 * time.Second)
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	command := m.Content
	fmt.Printf("监听的消息内容是：%s\n", command)
	switch command {
	case "!ping":
		data, err := motd.MotdBE(config.GetConfig().Host)
		if err != nil {
			fmt.Println("获取motd状态失败 错误：", err)
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, "服务器延迟为："+strconv.Itoa(int(data.Delay)))
	case "!cmd":
		fmt.Println("!cmd")
		return
	}
}
