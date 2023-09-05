package db

import (
	"encoding/json"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Store struct {
	//config *config.Config
	db *gorm.DB
}

// 结构体
type Wlist struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
	GameName string `gorm:"type:varchar(200); not null" json:"gameName" binding:"required"`
}

func (s *Store) init() {
	var err error
	host := config.GetConfig().Mysql.Host
	prot := config.GetConfig().Mysql.Port
	account := config.GetConfig().Mysql.Account
	password := config.GetConfig().Mysql.Password
	name := config.GetConfig().Mysql.Name

	dsn := account + ":" + password + "@tcp(" + host + ":" + prot + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"
	s.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.Error("mysql数据库连接失败:", err)
		return
	}

	sqlDB, err := s.db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 10 秒钟

	s.db.AutoMigrate(&Wlist{})
}

func Mysqladd(msg string) {
	s := &Store{}
	s.init()
	var data Wlist
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		logger.Error("mysql数据解析失败:", err)
		return
	}
	s.db.Create(&data)
}
