package danger

import (
	"github.com/gucooing/bdstobot/pkg/logger"
	"os/exec"
)

func Cmdstart(msg string) string {
	cmd := exec.Command("cmd.exe", "/C", msg)
	err := cmd.Start()
	if err != nil {
		logger.Warn().Msgf("执行失败:%d", err)
		return "执行失败"
	}
	err = cmd.Wait()
	if err != nil {
		logger.Warn().Msgf("执行失败:%d", err)
		return "执行失败"
	}
	return "执行成功"
}
