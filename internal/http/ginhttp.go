package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/db"
	"github.com/gucooing/bdstobot/pkg/logger"
	"net/http"
)

func Httpserver() {
	// 创建接口
	r := gin.Default()

	r.GET("/api/user/list", func(c *gin.Context) {
		userlist := db.SaveGameUserList()
		if userlist == "" {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
			return
		}
		c.JSON(200, userlist)
	})

	r.GET("/", func(c *gin.Context) {
		r.LoadHTMLFiles("data/html/index.html")
		r.NoRoute(gin.WrapH(http.FileServer(http.Dir("static"))))
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 设置端口号
	PORT := config.GetConfig().Port
	logger.Info("http服务启动在端口：%s 上", PORT)
	err := r.Run(":" + PORT)
	if err != nil {
		logger.Error("http服务启动失败:", err)
		return
	}
}
