package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"oauth2/config"
)

func ConnMysql() *gorm.DB{
	var err error
	cfg := config.Get();
	db, err := gorm.Open(cfg.Db.Default.Type, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Db.Default.User,
		cfg.Db.Default.Password,
		cfg.Db.Default.Host,
		cfg.Db.Default.Port,
		cfg.Db.Default.DbName))
	if err != nil {
		log.Fatalf("conn mysql fail:",err.Error())
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return db
}

