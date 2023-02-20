package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"net/url"
	"oauth2/config"
	"oauth2/model"
	"oauth2/utils"
	"oauth2/utils/session"
	"strings"
	"time"
)

type UserService struct {
	Db      *gorm.DB
	Srv     *server.Server
	Manager *manage.Manager
}

// AuthorizeHandler 授权请求处理
func (mysvc UserService) AuthorizeHandler(c *gin.Context) {
	var form url.Values
	if v, _ := session.Get(c.Request, "RequestForm"); v != nil {
		c.Request.ParseForm()
		if c.Request.Form.Get("client_id") == "" {
			form = v.(url.Values)
		}
	}

	c.Request.Form = form

	if err := session.Delete(c.Writer, c.Request, "RequestForm"); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := mysvc.Srv.HandleAuthorizeRequest(c.Writer, c.Request); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}

}

func (mysvc UserService) AuthHandler(c *gin.Context) {

	form, err := session.Get(c.Request, "LoggedInUserID")

	if err != nil {
		c.Writer.Header().Set("Location", "/login")
		c.Writer.WriteHeader(http.StatusFound)
		return
	}

	if form == nil {
		http.Error(c.Writer, "Invalid Request", http.StatusBadRequest)
		return
	}
	mt := model.MyClaims{
		UserType: []string{"user", "admin"},
		UserName: form.(string),
	}
	mysvc.Manager.MapAccessGenerate(&mt)

	form, err = session.Get(c.Request, "RequestForm")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if form == nil {
		http.Error(c.Writer, "Invalid Request", http.StatusBadRequest)
		return
	}
	clientID := form.(url.Values).Get("client_id")
	scope := form.(url.Values).Get("scope")
	data := TplData{
		Client: config.GetClient(clientID),
		Scope:  config.ScopeFilter(clientID, scope),
	}
	c.HTML(http.StatusOK, "auth.html", data)
}

type TplData struct {
	Client config.Client
	// 用户申请的合规scope
	Scope []config.Scope
	Error string
}

func (mysvc UserService) LoginHandler(c *gin.Context) {
	form, err := session.Get(c.Request, "RequestForm")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if form == nil {
		http.Error(c.Writer, "Invalid Request", http.StatusBadRequest)
		return
	}
	clientID := form.(url.Values).Get("client_id")
	scope := form.(url.Values).Get("scope")
	// 页面数据
	data := TplData{
		Client: config.GetClient(clientID),
		Scope:  config.ScopeFilter(clientID, scope),
	}
	if data.Scope == nil {
		http.Error(c.Writer, "Invalid Scope", http.StatusBadRequest)
		return
	}
	if c.Request.Method == "POST" {
		if c.Request.Form == nil {
			if err := c.Request.ParseForm(); err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		//账号密码验证
		if c.Request.Form.Get("type") == "password" {
			username := c.Request.Form.Get("username")
			password := c.Request.Form.Get("password")
			var User model.Oauth2User
			if err := mysvc.Db.Where("user_name=? and password= ?", username, password).Find(&User).Error; err != nil {
				data.Error = "用户名密码错误!"
				c.HTML(http.StatusOK, "login.html", data)
				return
			}
			mt := model.MyClaims{
				UserType: []string{"user", "admin"},
				UserName: username,
			}
			mysvc.Manager.MapAccessGenerate(&mt)
		}
		if err := session.Set(c.Writer, c.Request, "LoggedInUserID", c.Request.Form.Get("username")); err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		//c.Writer.Header().Set("Location", "/auth")
		//c.Writer.Header().Set("Location", "/hello")
		//c.Writer.WriteHeader(http.StatusFound)
		c.Redirect(http.StatusFound, "/auth")
		return
	}
	c.HTML(http.StatusOK, "login.html", data)
}

func (mysvc UserService) UserInfo(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 2003,
			"msg":  "请求头中auth为空",
		})
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusOK, gin.H{
			"code": 2004,
			"msg":  "请求头中auth格式有误",
		})
		c.Abort()
		return
	}
	// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
	mc, err := utils.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2005,
			"msg":  "无效的Token",
		})
		return
	}
	if mc.UserName == "" {
		mc.UserName = mc.Subject
	}
	// 将当前请求的username信息保存到请求的上下文c上
	c.JSON(http.StatusOK, gin.H{
		"userName": mc.UserName,
	})
	return
}

func (mysvc UserService) LogoutHandler(c *gin.Context) {
	if c.Request.Form == nil {
		if err := c.Request.ParseForm(); err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	redirectURI := c.Request.Form.Get("redirect_uri")
	//解码重定向地址
	enEscapeUrl, _ := url.QueryUnescape(redirectURI)
	if err := session.Delete(c.Writer, c.Request, "LoggedInUserID"); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("redirect_uri: %v\n", enEscapeUrl)
	c.Writer.Header().Set("Location", enEscapeUrl)
	c.Writer.WriteHeader(http.StatusFound)
}

func (mysvc UserService) TokenHandler(c *gin.Context) {
	err := mysvc.Srv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (mysvc UserService) TestHandler(c *gin.Context) {
	token, err := mysvc.Srv.ValidationBearerToken(c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	cli, err := mysvc.Manager.GetClient(c.Request.Context(), token.GetClientID())
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	apiname := c.Request.FormValue("apiname")
	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"user_id":    token.GetUserID(),
		"client_id":  token.GetClientID(),
		"scope":      token.GetScope(),
		"domain":     cli.GetDomain(),
		"apiname":    apiname,
	}

	//todo apiname和scope的校验 domain的校验 过期时间的校验
	e := json.NewEncoder(c.Writer)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func (mysvc UserService) LoginHandlerTEST(c *gin.Context) {
	form, err := session.Get(c.Request, "RequestForm")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if form == nil {
		http.Error(c.Writer, "Invalid Request", http.StatusBadRequest)
		return
	}
	clientID := form.(url.Values).Get("client_id")
	scope := form.(url.Values).Get("scope")

	// 页面数据
	data := TplData{
		Client: config.GetClient(clientID),
		Scope:  config.ScopeFilter(clientID, scope),
	}
	if data.Scope == nil {
		http.Error(c.Writer, "Invalid Scope", http.StatusBadRequest)
		return
	}

	if c.Request.Method == "POST" {
		if c.Request.Form == nil {
			if err := c.Request.ParseForm(); err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		//账号密码验证
		/*		if r.Form.Get("type") == "password" {
				//自己实现验证逻辑
				var user model.User
				userID := user.GetUserIDByPwd(r.Form.Get("username"), r.Form.Get("password"))
				if userID == "" {
					t, err := tpl.ParseFiles("tpl/login.html")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					data.Error = "用户名密码错误!"
					t.Execute(w, data)

					return
				}
			}*/
		/*		if r.Form.Get("type") == "password" {
				fmt.Println("进行密码校验 ")
				//自己实现验证逻辑
				username := r.Form.Get("username")
				password := r.Form.Get("password")
				err := ExampleSearch(username, password)
				if err != nil {
					t, err := tpl.ParseFiles("tpl/login.html")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					data.Error = "用户名密码错误!"
					t.Execute(w, data)
					return
				}*/

		//登陆成功,生成自定义的jwt
		// mt := model.MyClaims{
		// 	UserType: []string{"user", "admin"},
		// 	UserName: "daijeizou",
		// }
		// mysvc.Manager.MapAccessGenerate(&mt)

	}
	//扫码验证
	//手机验证码验证

	//登陆成功返回access_token
	//manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))
	//登陆成功设置session值
	if err := session.Set(c.Writer, c.Request, "LoggedInUserID", c.Request.Form.Get("username")); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Location", "/auth")
	c.Writer.WriteHeader(http.StatusFound)
	return
}
