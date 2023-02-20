package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

func MakePwd(password string)(salt string,newpwd string)  {
	salt = base64.StdEncoding.EncodeToString([]byte(time.Now().String()))
	fmt.Println(salt)
	m5 := md5.New()
	m5.Write([]byte(password))
	m5.Write([]byte(salt))
	st := m5.Sum(nil)
	newpwd = hex.EncodeToString(st)
	return salt,newpwd
}

func CheckPwd(salt string,password string,enCodepwd string)bool   {
	m5 := md5.New()
	m5.Write([]byte(password))
	m5.Write([]byte(salt))
	st := m5.Sum(nil)
	Inputpwd := hex.EncodeToString(st)
	return Inputpwd == enCodepwd
}