package http

func (s *Server) initRouter() {
	s.router.Use(clientIPMiddleware())
	s.router.Any("/", s.handleDefault)
	s.router.Any("/index.html", s.handleDefault)
	s.router.GET("/api/user/list", s.handleApiListand)          //查询全部用户数据
	s.router.GET("/api/user/add", s.hadnleApiUseradd)           //增加用户数据
	s.router.GET("/api/user/delete/:id", s.hadnleApiUserdelete) //删除用户数据
	s.router.GET("/api/user/update/:id", s.hadnleApiUserupdate) //更新用户数据
}
