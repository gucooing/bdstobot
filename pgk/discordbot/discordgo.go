package discordbot

import (
	"encoding/json"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/motd"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	GuildID         string `json:"GuildID"`
	DiscordBotToken string `json:"DiscordBotToken"`
}

var s *discordgo.Session

func init() {
	var err error
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("无法读取配置文件: %v\n", err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("配置文件解析错误: %v\n", err)
	}
	s, err = discordgo.New("Bot " + config.DiscordBotToken)
	if err != nil {
		fmt.Println("discord bot token 无效: %v", err)
	}
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "使用ping 测试服务器延迟",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseCN: "ping",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseCN: "使用ping 测试服务器延迟",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "选项",
					Description: "此指令的一级选项",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseCN: "选项",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseCN: "此指令的一级选项",
					},
					Type: discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "选项参数一",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseCN: "选项参数一",
							},
							Value: 1,
						},
						{
							Name: "选项参数二",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseCN: "选项参数er",
							},
							Value: 2,
						},
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data, _ := motd.MotdBE(config.GetConfig().Host)
			responses := map[discordgo.Locale]string{
				discordgo.ChineseCN: "服务器延迟为：" + strconv.Itoa(int(data.Delay)),
			}
			response := "服务器延迟为：" + strconv.Itoa(int(data.Delay))
			if r, ok := responses[i.Locale]; ok {
				response = r
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func DiscordBot() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("登录bot: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("bot无法连接到discord: %v\n", err)
	}
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("无法读取配置文件: %v\n", err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("配置文件解析错误: %v\n", err)
	}

	fmt.Println("注册命令中...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, config.GuildID, v)
		if err != nil {
			log.Fatalf("无法注册 '%v' 命令: %v\n", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	fmt.Println("discord bot 命令已成功注册 !")
	for {
		// 阻塞携程，保持机器人在线
		time.Sleep(10 * time.Second)
	}
}
