package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	LogLevel string `json:"logLevel"`
	Port     string `json:"port"`
	McHost   string `json:"McHost"`
	Mcpath   string `json:"Mcpath"`
	Pflp     *pflp
	QQ       *qq
	Discord  *discord
	ZmHost   string `json:"ZmHost"`
	Mysql    *mysql
}

type pflp struct {
	PFLPWsurl string `json:"PFLPWsurl"`
	Key       string `json:"Key"`
}

type discord struct {
	DiscordWebhookUrl string `json:"DiscordWebhookUrl"`
	DiscordBot        bool
	DiscordWsurl      string `json:"DiscordWsurl"`
	DiscordBotToken   string `json:"DiscordBotToken"`
	GuildID           string `json:"GuildID"`
}

type qq struct {
	QQ          bool
	QqAdmin     int64  `json:"Qqadmin"`
	QQgroup     int64  `json:"QQgroup"`
	CqhttpWsurl string `json:"CqhttpWsurl"`
}

type mysql struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Password string `json:"password"`
	BdTable  string `json:"bd_table"`
}

var CONF *Config = nil

func GetConfig() *Config {
	return CONF
}

var FileNotExist = errors.New("config file not found")

func LoadConfig() error {
	filePath := "./config.json"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	f, err := os.Open(filePath)
	if err != nil {
		return FileNotExist
	}
	defer func() {
		_ = f.Close()
	}()
	c := new(Config)
	d := json.NewDecoder(f)
	if err := d.Decode(c); err != nil {
		return err
	}
	CONF = c
	return nil
}

var DefaultConfig = &Config{
	LogLevel: "Info",
	Port:     "8080",
	McHost:   "127.0.0.1:19132",
	Mcpath:   "D:\\bedrock_server.exe",
	Pflp: &pflp{
		PFLPWsurl: "ws://127.0.0.1:80",
		Key:       "1234567890",
	},
	QQ: &qq{
		QQ:          false,
		QqAdmin:     123456789,
		QQgroup:     123456789,
		CqhttpWsurl: "ws://127.0.0.1:80",
	},
	Discord: &discord{
		DiscordWebhookUrl: "https://127.0.0.1",
		DiscordBot:        false,
		DiscordWsurl:      "ws://127.0.0.1:80",
		DiscordBotToken:   "1234567890",
		GuildID:           "",
	},
	ZmHost: "127.0.0.1:16261",
	Mysql: &mysql{
		Host:     "127.0.0.1",
		Port:     "3306",
		Name:     "root",
		Password: "123456789",
		BdTable:  "list",
	},
}
