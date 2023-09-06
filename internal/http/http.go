package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/internal/db"
	"github.com/gucooing/bdstobot/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var (
	// mutex用于确保同时只有一个请求可以触发重启
	mutex sync.Mutex
)

type Server struct {
	config *config.Config
	Store  *db.Store
	router *gin.Engine
	server *http.Server
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{}
	s.config = cfg
	s.Store = db.NewStore(s.config) // 初始化数据库
	if s.Store == nil {
		fmt.Println("初始化数据库失败")
		return nil
	}

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.Use(gin.Recovery())
	s.router.LoadHTMLFiles("data/html/index.html")
	s.router.NoRoute(gin.WrapH(http.FileServer(http.Dir("static"))))

	return s
}

func (s *Server) Start() error {
	// 初始化路由
	s.initRouter()

	// 获取地址
	addr := s.config.Addr

	go s.startServer(addr, "HTTP")

	return nil
}

// startServer 启动一个 HTTP 服务器。
func (s *Server) startServer(addr string, serverType string) {
	s.server = &http.Server{Addr: addr, Handler: s.router}
	logger.Info("listen_addr: %s, %s服务器正在启动", addr, serverType)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Info("listen_addr: %s, %s服务器正在启动", addr, serverType)

	}
}

func Restart() error {
	mutex.Lock()
	defer mutex.Unlock()
	// 获取当前进程的可执行文件路径
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get the executable path: %w", err)
	}
	// 创建一个Cmd结构体，用于启动新进程
	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	// 开始新进程
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start new process: %w", err)
	}
	// 结束当前进程
	os.Exit(0)
	return nil
}
func (s *Server) Shutdown(context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Close()
}

func clientIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			c.Next()
			return
		}

		// 将 IP 信息存储在 gin.Context 中
		c.Set("IP", ip)

		c.Next()
	}
}
