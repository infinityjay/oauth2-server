package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"net/http"
	"oauth2/config"
	_ "oauth2/docs"
	"oauth2/handlers"
	"oauth2/handlers/Oauth2Config"
	"oauth2/model"
	"oauth2/utils"
	"oauth2/utils/session"
)

func main() {
	config.Setup()
	// init db connection
	// configure db in app.yaml then uncomment
	// model.Setup()
	session.Setup()

	// manager config
	manager := manage.NewDefaultManager()
	//设置一些token的有效时间
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// or use redis token store
	//manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
	//    Addr: config.Get().Redis.Default.Addr,
	//    DB: config.Get().Redis.Default.Db,
	//}))
	jwt := new(model.MyClaims)
	manager.MapAccessGenerate(jwt)

	//保存客户端的的相关信息
	clientStore := store.NewClientStore()
	for _, v := range config.Get().OAuth2.Client {
		clientStore.Set(v.ID, &models.Client{
			ID:     v.ID,
			Secret: v.Secret,
			Domain: v.Domain,
		})
	}
	manager.MapClientStorage(clientStore)
	// config oauth2 server
	srv := server.NewDefaultServer(manager)

	//对服务器进行设置，使其能够进行基于账号密码登陆
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)

	//在此处设置了登陆跳转的地址
	srv.SetUserAuthorizationHandler(Oauth2Config.UserAuthorizeHandler)
	//设置允许的授权模式类型
	//srv.SetAllowedResponseType(oauth2.Token)

	//设置自定义的权限范围
	srv.SetAuthorizeScopeHandler(Oauth2Config.AuthorizeScopeHandler)
	srv.SetInternalErrorHandler(Oauth2Config.InternalErrorHandler)
	srv.SetResponseErrorHandler(Oauth2Config.ResponseErrorHandler)
	myUserSvc := handlers.UserService{
		Srv:     srv,
		Manager: manager,
		Db:      utils.ConnMysql(),
	}
	r := gin.Default()
	r.LoadHTMLGlob("tpl/*")
	r.StaticFS("/static", http.Dir("./static")) //配置静态文件夹路径

	//swagger api
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	//授权服务
	r.Any("/authorize", myUserSvc.AuthorizeHandler)
	r.GET("/auth", myUserSvc.AuthHandler)
	r.GET("/user", myUserSvc.UserInfo)
	r.Any("/login", myUserSvc.LoginHandler)
	r.GET("/logout", myUserSvc.LogoutHandler)
	r.POST("/token", myUserSvc.TokenHandler)
	r.GET("/test", myUserSvc.TestHandler)

	log.Println("Server is running at 9096 port.")
	r.Run(":9096")
}
func passwordAuthorizationHandler(username, password string) (userID string, err error) {
	var user model.User
	userID = user.GetUserIDByPwd(username, password)

	return
}
