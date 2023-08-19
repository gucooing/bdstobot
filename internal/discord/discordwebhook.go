package discord

import (
	"bytes"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/logger"
	"io/ioutil"
	"net/http"
)

func Discordwebhook(msg string) {
	url := config.GetConfig().DiscordWebhookUrl
	body := "{\"content\": \"" + msg + "\", \"username\": \"MCBDS\", \"avatar_url\": \"https://webusstatic.yo-star.com/bluearchive_jp_web/fankit/162704158443017840/01.png\"}"
	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		logger.Warn().Msgf("discordbot webhook消息发送错误:%d", err)
	} else {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Warn().Msgf("discordbot webhook接收回调body失败:%d", err)
			return
		}
		logger.Debug().Msgf("discordbot webhook回调body：%d", string(b))
		return
	}
	return
}
