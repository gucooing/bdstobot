package danger

import (
	"fmt"
	"os/exec"
)

func Cmdstart(msg string) string {
	cmd := exec.Command("cmd.exe", "/C", msg)
	err := cmd.Start()
	if err != nil {
		fmt.Println("执行失败:", err)
		return "执行失败"
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("执行失败:", err)
		return "执行失败"
	}
	return "执行成功"
}
