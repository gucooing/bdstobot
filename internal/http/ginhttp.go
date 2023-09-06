package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/db"
	"github.com/gucooing/bdstobot/pkg/logger"
	"net/http"
	"strconv"
)

var (
	s    db.Store
	data db.Wlist
)

func Httpserver() {
	// 创建接口
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		r.LoadHTMLFiles("data/html/index.html")
		r.NoRoute(gin.WrapH(http.FileServer(http.Dir("static"))))
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	//查全部
	r.GET("/api/user/list", func(c *gin.Context) {
		// 1. 查询全部数据,  查询分页数据
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))
		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}
		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}
		// 返回一个总数
		var total int64
		// 查询数据库
		datas := s.Mysqllistand(total, pageSize, offsetVal, &data)
		if len(datas) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     datas,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})

	//增
	r.GET("/api/user/add", func(c *gin.Context) {
		err := c.ShouldBindJSON(&data)
		// 判断绑定是否有错误
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			// 数据库的操作
			s.Mysqladd(&data)

			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})

	// 删
	r.GET("/api/user/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		boold := s.Mysqldelete(id)
		if boold {
			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "id没有找到, 删除失败",
				"code": 400,
			})
		}
	})

	//修改
	r.GET("/api/user/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		// 判断 id 是否存在
		boolif := s.Mysqllist(id)
		if boolif {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				// db 修改数据库内容
				s.Mysqlupdate(id, &data)
				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		} else {
			c.JSON(200, gin.H{
				"msg":  "用户id没有找到",
				"code": 400,
			})
		}
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
