package utils

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"log"
)

func ExampleSearch(username string, password string) error {
	l, err := ldap.DialURL("ldap://119.3.31.108:389")
	if err != nil {
		log.Fatal("连接ldap服务失败:", err)
	}
	defer l.Close()

	//后续改成配置
	_, err = l.SimpleBind(&ldap.SimpleBindRequest{
		Username: "cn=admin,dc=youedata,dc=com",
		Password: "youedata520",
	})
	if err != nil {
		log.Fatalf("Failed to bind: %s\n", err)
	}
	searchRequest := ldap.NewSearchRequest(
		"dc=youedata,dc=com", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		//"(&(objectClass=organizationalPerson))", // 查询所有人
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username), //查询指定人
		[]string{"dn", "cn"},                                                   // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) == 0 {
		return fmt.Errorf("无效的用户名或id")
	}
	fmt.Println(sr.Entries[0].Attributes[0].Values)
	fmt.Println(sr.Entries[0].Attributes[0].Name)

	//校验用户密码
	err = l.Bind(sr.Entries[0].DN, password)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	return nil

}
