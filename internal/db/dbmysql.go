package db

import (
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var data []Wlist

type Store struct {
	config *config.Config
	Db     *gorm.DB
}

// 结构体
type Wlist struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
	GameName string `gorm:"type:varchar(200); not null" json:"gameName" binding:"required"`
}

func (s *Store) init() {
	var err error
	host := s.config.Mysql.Host
	prot := s.config.Mysql.Port
	account := s.config.Mysql.Account
	password := s.config.Mysql.Password
	name := s.config.Mysql.Name

	dsn := account + ":" + password + "@tcp(" + host + ":" + prot + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"
	s.Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.Error("mysql数据库连接失败:", err)
		return
	}
	sqlDB, err := s.Db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 10 秒钟
	//defer sqlDB.Close()

	s.Db.AutoMigrate(&Wlist{})
}

// NewStore 创建一个新的 store。
func NewStore(config *config.Config) *Store {
	s := &Store{config: config}
	s.init()
	return s
}

// 查全部
func (s *Store) Mysqllistand(total int64, pageSize, offsetVal int, datal *Wlist) []Wlist {
	s.Db.Model(datal).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&data)
	if len(data) == 0 {
		return nil
	} else {
		return data
	}
}

// 查一个
func (s *Store) Mysqllist(id string) bool {
	s.Db.Where("id = ?", id).First(&data)
	if len(data) == 0 {
		return false
	}
	return true
}

// 增加
func (s *Store) Mysqladd(data *Wlist) {
	s.Db.Create(&data)
}

// 删
func (s *Store) Mysqldelete(id string) bool {
	s.Db.Where("id = ?", id).First(&data)
	if len(data) == 0 {
		return false
	}
	s.Db.Where("id = ?", id).Delete(&data)
	return true
}

// 改
func (s *Store) Mysqlupdate(id string, data *Wlist) {
	s.Db.Where("id = ?", id).Updates(&data)
}
