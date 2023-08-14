package pgk

import (
	"bytes"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"io/ioutil"
	"net/http"
)

func Discord(msg string) {
	url := config.GetConfig().DiscordWebhookUrl
	body := "{\"content\": \"" + msg + "\", \"username\": \"MCBDS\", \"avatar_url\": \"https://webusstatic.yo-star.com/bluearchive_jp_web/fankit/162704158443017840/01.png\"}"
	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("discord webhook消息发送错误:", err)
	}
	b, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(b))
}
