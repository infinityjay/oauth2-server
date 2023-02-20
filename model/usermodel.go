package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type BaseUserModel struct {
	ID         int       `gorm:"column:Id;primary_key;form:id" `
	Username   string    `gorm:"column:Username" form:"username"`
	Phone      string    `gorm:"column:Phone" form:"phone"`
	Status     int       `gorm:"column:Status" form:"status"`
	CreateTime time.Time `gorm:"column:CreateTime" form:"createTime"`
	UpdateTime time.Time `gorm:"column:UpdateTime" form:"updateTime"`
}

type UserModel struct {
	BaseUserModel
	Password           string `gorm:"column:Password" form:"password"`
	ClientScopeType    int    `gorm:"column:ClientScopeType"`
	NeedChangePassword int    `gorm:"column:NeedChangePassword" form:"needChangePassword"`
	Salt               string `gorm:"column:Salt"`
}

type Oauth2User struct {
	UserId   int    `json:"user_id" xorm:"pk autoincr BIGINT"`
	UserName string `json:"user_name" xorm:"VARCHAR(256)"`
	//Password     string `json:"password" xorm:"VARCHAR(256)`
	Salt         string `json:"salt" xorm:"VARCHAR(256)"`
	DisplayName  string `json:"display_name" xorm:"VARCHAR(256)"`
	Phone        string `json:"phone" xorm:"VARCHAR(256)"`
	Email        string `json:"email" xorm:"VARCHAR(256)"`
	Status       int    `json:"status" xorm:"not null default 1 TINYINT"` // 0=正常;1=禁用
	Source       string `json:"source" xorm:"VARCHAR(256)"`
	CreateUserId int    `json:"create_user_id" xorm:"BIGINT"`
	UpdateUserId int    `json:"update_user_id" xorm:"BIGINT"`
	CreateTime   string `json:"create_time" xorm:"DATETIME"`
	UpdateTime   string `json:"update_time" xorm:"DATETIME"`
	IsDelete     int    `json:"is_delete" xorm:"not null default 0 TINYINT"` // 0=正常;1=删除
}

func (u *Oauth2User) TableName() string {
	return "sys_user"
}

type CreateUserModel struct {
	UserModel
	UserRoles []string `form:"userRoles" gorm:"-"`
}

type ViewUserModel struct {
	BaseUserModel
	UserRoles []string `form:"userRoles" sql:"-"`
}

type Role struct {
	Id   string `gorm:"column:ID`
	Name string `gorm:"column:Name"`
}

type UserRole struct {
	//UserModel UserModel `gorm:"ForeignKey:UserId;-"`
	UserId int `gorm:"column:UserId"`
	//Role      Role      `gorm:"ForeignKey:RoleIdId;-"`
	RoleId string `gorm:"column:RoleId"`
}

func (u *Role) TableName() string {
	return "role"
}

func (u *UserRole) TableName() string {
	return "user_role"
}

func (u *BaseUserModel) TableName() string {
	return "user"
}
