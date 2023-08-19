package db

import (
	"encoding/json"
	"github.com/gucooing/bdstobot/pkg/logger"
	"os"
)

type GameInfo struct {
	QQ       int64  `json:"qq"`
	GameName string `json:"gamename"`
}

func SaveGameInfo(qq int64, gamename string) bool {
	// 读取已有的游戏信息
	gameInfos, err := LoadGameInfos()
	if err != nil {
		logger.Warn("读取 JSON 文件中游戏信息失败:", err)
		// 如果读取失败，则创建一个新的切片用于存储游戏信息
		gameInfos = make([]GameInfo, 0)
	}
	// 添加新的游戏信息到切片中
	gameInfos = append(gameInfos, GameInfo{
		QQ:       qq,
		GameName: gamename,
	})
	// 将游戏信息保存到 JSON 文件中
	bytes, err := json.Marshal(gameInfos)
	if err != nil {
		logger.Warn("将游戏信息保存到 JSON 文件中失败:", err)
		return false
	}
	err = os.WriteFile("data/game.json", bytes, 0644)
	if err != nil {
		logger.Warn("写入 JSON 文件中失败:", err)
		return false
	}
	return true
}

func DeleteGameInfoByQQ(qq int64) bool {
	// 读取游戏信息
	gameInfos, err := LoadGameInfos()
	if err != nil {
		logger.Warn("读取 JSON 文件中游戏信息失败:", err)
		return false
	}

	// 查找并删除匹配的 QQ 号码的游戏信息
	found := false
	for i, info := range gameInfos {
		if info.QQ == qq {
			// 删除匹配的游戏信息
			gameInfos = append(gameInfos[:i], gameInfos[i+1:]...)
			found = true
			break
		}
	}
	// 如果找到匹配的 QQ 号码，则保存更新后的游戏信息到 JSON 文件中
	if found {
		bytes, err := json.Marshal(gameInfos)
		if err != nil {
			logger.Warn("将游戏信息保存到 JSON 文件中失败:", err)
			return false
		}
		err = os.WriteFile("data/game.json", bytes, 0644)
		if err != nil {
			logger.Warn("写入 JSON 文件中失败:", err)
			return false
		}
	} else {
		return false
	}
	return true
}

func FindGameNameByQQ(qq int64) string {
	// 读取游戏信息
	gameInfos, err := LoadGameInfos()
	if err != nil {
		logger.Warn("读取 JSON 文件中游戏信息失败:", err)
		return ""
	}

	// 遍历游戏信息，查找匹配的 QQ 号码
	for _, info := range gameInfos {
		if info.QQ == qq {
			return info.GameName
		}
	}

	// 如果没有找到匹配的 QQ 号码，则返回空字符串和 false
	return ""
}

func LoadGameInfos() ([]GameInfo, error) {
	// 从 JSON 文件中读取游戏信息
	bytes, err := os.ReadFile("data/game.json")
	if err != nil {
		logger.Warn("读取 JSON 文件中游戏信息失败:", err)
		return nil, err
	}
	// 解析 JSON 数据到切片中
	var gameInfos []GameInfo
	err = json.Unmarshal(bytes, &gameInfos)
	if err != nil {
		logger.Warn("解析 JSON 文件中游戏信息失败:", err)
		return nil, err
	}
	return gameInfos, nil
}
