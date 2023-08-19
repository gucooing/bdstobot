package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	LogLevel          string `json:"logLevel"`
	Host              string `json:"Host"`
	QQ                bool
	QqAdmin           int64  `json:"Qqadmin"`
	QQgroup           int64  `json:"QQgroup"`
	CqhttpWsurl       string `json:"CqhttpWsurl"`
	PFLPWsurl         string `json:"PFLPWsurl"`
	DiscordWebhookUrl string `json:"DiscordWebhookUrl"`
	DiscordBot        bool
	DiscordWsurl      string `json:"DiscordWsurl"`
	DiscordBotToken   string `json:"DiscordBotToken"`
	GuildID           string `json:"GuildID"`
	Key               string `json:"Key"`
	Mcpath            string `json:"mcpath"`
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
	LogLevel:          "Info",
	Host:              "127.0.0.1:19132",
	QQ:                false,
	QqAdmin:           123456789,
	QQgroup:           123456789,
	CqhttpWsurl:       "ws://127.0.0.1:80",
	PFLPWsurl:         "ws://127.0.0.1:80",
	DiscordWebhookUrl: "https://127.0.0.1",
	DiscordBot:        false,
	DiscordWsurl:      "ws://127.0.0.1:80",
	DiscordBotToken:   "1234567890",
	GuildID:           "",
	Key:               "1234567890",
	Mcpath:            "D:\\bedrock_server.exe",
}
