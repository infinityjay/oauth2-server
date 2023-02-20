package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"oauth2/model"
	"oauth2/utils"
	"strconv"
	"time"
)

func (mysvc UserService) AddUser(c *gin.Context) {
	var createuser model.CreateUserModel
	//获取表单数据
	if err := c.Bind(&createuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createuser.UpdateTime = time.Now()
	createuser.CreateTime = time.Now()
	salt,newpwd := utils.MakePwd(createuser.Password)
	createuser.Password = newpwd
	createuser.Salt = salt

	//开始事务
	tx := mysvc.Db.Begin()
	//创建用户
	if err := tx.Create(&createuser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("id:",createuser.ID)
	for _, role := range createuser.UserRoles {
		userrole := model.UserRole{
			RoleId: role,
			UserId: createuser.ID,
		}
		if err := tx.Create(&userrole).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  2001,
				"error": err.Error(),
			})
			return
		}

	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})

}



func (mysvc UserService) GetUserInfo(c *gin.Context) {
	ID := c.Param("id")
	var user model.BaseUserModel
	if err := mysvc.Db.Where("ID=?", ID).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": err.Error(),
		})
		return
	}
	var userRoles []model.UserRole
	if err := mysvc.Db.Where("UserId=?", ID).Find(&userRoles).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": err.Error(),
		})
		return
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.RoleId)
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"phone":      user.Phone,
		"status":     user.Status,
		"createTime": user.CreateTime,
		"updateTime": user.UpdateTime,
		"userRoles":  roles,
	})
}

func (mysvc UserService) GetRoles(c *gin.Context)  {
	var roles []model.Role
	if err:=mysvc.Db.Find(&roles).Error;err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,roles)

}

func (mysvc UserService) GetUsers(c *gin.Context) {
	PageNo, _ := strconv.Atoi(c.Query("PageNo"))
	PageSize, _ := strconv.Atoi(c.Query("PageSize"))
	var PageCount int
	var total int
	mysvc.Db.Table("user").Count(&total) //计算表中所有用户的数量
	var users []model.BaseUserModel
	mysvc.Db.Limit(PageSize).Offset((PageNo - 1) * PageSize).Find(&users)
	PageCount = len(users)
	var viewUsers []model.ViewUserModel
	for _, user := range users {
		var userRoles []model.UserRole
		var roles []string
		if err := mysvc.Db.Where("UserId=?", user.ID).Find(&userRoles).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg": err.Error(),
			})
			return
		}
		for _, role := range userRoles {
			roles = append(roles, role.RoleId)
		}
		viewuser := model.ViewUserModel{
			user,
			roles,
		}
		viewUsers = append(viewUsers, viewuser)
	}
	c.JSON(http.StatusOK, gin.H{
		"items":     viewUsers,
		"pageNo":    PageNo,
		"pageCount": PageCount,
		"pageSize":  PageSize,
		"total":     total,
	})
}

func (mysvc UserService) ChangeUser(c *gin.Context) {
	id,_ := strconv.Atoi(c.Param("id"))

	var createuser model.CreateUserModel
	//获取表单数据
	if err := c.Bind(&createuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createuser.UpdateTime = time.Now()
	salt,newpwd := utils.MakePwd(createuser.Password)
	tx := mysvc.Db.Begin()
	if err := tx.Model(&createuser).Where("ID=?", id).Updates(map[string]interface{}{
		"Username": createuser.Username, "Password":newpwd, "Phone": createuser.Phone, "Status": createuser.Status,"Salt":salt}).Error;err!=nil{
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := tx.Where("UserId=?",id).Delete(model.UserRole{}).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, role := range createuser.UserRoles {
		userrole := model.UserRole{
			RoleId: role,
			UserId: id,
		}
		if err := tx.Create(&userrole).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  2001,
				"error": err.Error(),
			})
			return
		}

	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})

}

func (mysvc UserService) Change(c *gin.Context) {
	var user model.UserModel
	//获取表单数据
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//创建用户
	if err := mysvc.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
