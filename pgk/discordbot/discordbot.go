package discordbot

import (
	"encoding/json"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pgk/motd"
	"github.com/gucooing/bdstobot/pgk/takeover"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	GuildID         string `json:"GuildID"`
	DiscordBotToken string `json:"DiscordBotToken"`
	RemoveCommands  bool   `json:"RemoveCommands"`
}

var s *discordgo.Session

func init() {
	var err error
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("无法读取配置文件: %v\n", err)
	}
	var nweconfig Config
	err = json.Unmarshal(file, &nweconfig)
	if err != nil {
		fmt.Printf("配置文件解析错误: %v\n", err)
	}
	s, err = discordgo.New("Bot " + nweconfig.DiscordBotToken)
	if err != nil {
		fmt.Printf("discord bot token 无效: %v\n", err)
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
		},
		{
			Name:        "绑定",
			Description: "使用“绑定”指令添加服务器白名单",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseCN: "绑定",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseCN: "使用“绑定”指令添加服务器白名单",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "绑定",
					Description: "你游戏里面的昵称",
					Required:    true,
				},
			},
		},
		{
			Name:        "解绑",
			Description: "使用“解绑”指令删除服务器白名单",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseCN: "解绑",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseCN: "使用“解绑”指令删除服务器白名单",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "解绑",
					Description: "你游戏里面的昵称",
					Required:    true,
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
				fmt.Println(err)
				return
			}
		},
		"绑定": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			msgformat := "操作成功:\n"

			user := i.Interaction.Member.User
			username := user.Username

			if option, ok := optionMap["绑定"]; ok {
				margs = append(margs, username, option.StringValue())
				//建议在此进行逻辑处理
				margss := "whitelist add " + option.StringValue()
				takeover.Pflpwsreq("cmd", margss)
				msgformat += "> 用户: %s\n> 游戏昵称: %s\n"
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		},
		"解绑": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			msgformat := "操作成功:\n"

			user := i.Interaction.Member.User
			username := user.Username

			if option, ok := optionMap["解绑"]; ok {
				margs = append(margs, username, option.StringValue())
				//建议在此进行逻辑处理
				margss := "whitelist remove " + option.StringValue()
				takeover.Pflpwsreq("cmd", margss)
				msgformat += "> 用户: %s\n> 游戏昵称: %s\n"
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
			if err != nil {
				fmt.Println(err)
				return
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
		fmt.Printf("bot无法连接到discord: %v\n", err)
		return
	}
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("无法读取配置文件: %v\n", err)
		return
	}
	var nweconfig Config
	err = json.Unmarshal(file, &nweconfig)
	if err != nil {
		fmt.Printf("配置文件解析错误: %v\n", err)
		return
	}

	fmt.Println("注册命令中...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, nweconfig.GuildID, v)
		if err != nil {
			fmt.Printf("无法注册 '%v' 命令: %v\n", v.Name, err)
			return
		}
		registeredCommands[i] = cmd
	}
	fmt.Println("discord bot 命令已成功注册 !")
	for {
		// 阻塞携程，保持机器人在线
		time.Sleep(10 * time.Second)
	}
}
