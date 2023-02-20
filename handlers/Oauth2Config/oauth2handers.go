package Oauth2Config

import (
	"github.com/go-oauth2/oauth2/v4/errors"
	"log"
	"net/http"
	"oauth2/config"
	"oauth2/utils/session"
)

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	//校验用户是否已经登陆
	v, _ := session.Get(r, "LoggedInUserID")
	if v == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		session.Set(w, r, "RequestForm", r.Form)
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = v.(string)

	// 不记住用户
	//session.Delete(w,r,"LoggedInUserID")
	// store.Delete("LoggedInUserID")
	// store.Save()
	return
}

// AuthorizeScopeHandler 根据client注册的scope 过滤非法scope
func AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	s := config.ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
	if s == nil {
		http.Error(w, "Invalid Scope", http.StatusBadRequest)
		return
	}
	scope = config.ScopeJoin(s)

	return
}

func InternalErrorHandler(err error) (re *errors.Response) {
	log.Println("Internal Error:", err.Error())
	return
}

func ResponseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}

