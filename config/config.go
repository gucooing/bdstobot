package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Host              string `json:"Host"`
	QqAdmin           int64  `json:"Qqadmin"`
	QQgroup           int64  `json:"QQgroup"`
	CqhttpWsurl       string `json:"CqhttpWsurl"`
	BdsWsurl          string `json:"BdsWsurl"`
	DiscordWebhookUrl string `json:"DiscordWebhookUrl"`
	DiscordBotToken   string `json:"DiscordBotToken"`
	GuildID           string `json:"GuildID"`
	Key               string `json:"Key"`
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
	Host:              "127.0.0.1:19132",
	QqAdmin:           123456789,
	QQgroup:           568052501,
	CqhttpWsurl:       "ws://127.0.0.1:80",
	BdsWsurl:          "http://127.0.0.1",
	DiscordWebhookUrl: "https://127.0.0.1",
	DiscordBotToken:   "1234567890",
	GuildID:           "",
	Key:               "1234567890",
}
