package discordbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type SendMessageInput struct {
	ChannelID string
	Content   string
}

func SendMessage(s *discordgo.Session, input *SendMessageInput) error {
	gs, err := s.ChannelMessageSend(input.ChannelID, input.Content)
	fmt.Printf("发送的地址是：%s\n , 发送的消息内容是：%s\n", input.ChannelID, input.Content)
	if err != nil {
		fmt.Println("discord bot 发送错误：", err)
		return err
	}

	fmt.Println("discord bot 成功发送：", gs)
	return nil
}
