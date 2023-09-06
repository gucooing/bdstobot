package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/bdstobot/internal/db"
	"net/http"
	"strconv"
)

var data db.Wlist

func (s *Server) handleDefault(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (s *Server) handleApiListand(c *gin.Context) {
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
	datas := s.Store.Mysqllistand(total, pageSize, offsetVal, &data)
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
}

func (s *Server) hadnleApiUseradd(c *gin.Context) {
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
		s.Store.Mysqladd(&data)

		c.JSON(200, gin.H{
			"msg":  "添加成功",
			"data": data,
			"code": 200,
		})
	}
}

func (s *Server) hadnleApiUserdelete(c *gin.Context) {
	id := c.Param("id")
	boold := s.Store.Mysqldelete(id)
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
}

func (s *Server) hadnleApiUserupdate(c *gin.Context) {
	id := c.Param("id")
	// 判断 id 是否存在
	boolif := s.Store.Mysqllist(id)
	if boolif {
		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "修改失败",
				"code": 400,
			})
		} else {
			// db 修改数据库内容
			s.Store.Mysqlupdate(id, &data)
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
}
