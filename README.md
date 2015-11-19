# gin-session
go gin框架的session中间件， 目前只实现了redis驱动

使用方法
func main(){
	g := gin.New()

	session_config := `{
		"storeDriver":"redis",
		"cookieName":"session_id",
		"EnableSetCookie":true,
		"secure":false,
		"cookieLifeTime": 3,
		"Domain":""
	}`
	s := g.Group("/")
	s.Use(gsession.Middleware(session_config))
	{
		s.GET("/",func(c *gin.Context){
			sess := gsession.GetSession(c);
			sess.Set("uid","22")
		})
		s.GET("/member",func(c *gin.Context){
			sess := gsession.GetSession(c);
			uid := sess.Get("uid")
			fmt.Println(uid)
		})
	}
	g.Run(":8888")
}




