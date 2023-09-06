package internal

import (
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/db"
)

type Server struct {
	Config *config.Config
	Store  *db.Store
}

func NewMysql(cfg *config.Config) *Server {
	s := &Server{}
	s.Config = cfg
	s.Store = db.Mysqlmain(s.Config) // 初始化数据库
	if s.Store == nil {
		fmt.Println("初始化数据库失败")
		return nil
	}

	return s
}
