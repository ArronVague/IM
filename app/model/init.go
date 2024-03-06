package model

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var DbEngine *xorm.Engine

func init() {
	fmt.Println("init database")
	driverName := "mysql"
	dsnName := "root:Wohaijidemima123@(gz-cynosdbmysql-grp-04b3z61j.sql.tencentcdb.com:20464)/chat?charset=utf8"
	err := errors.New("")
	DbEngine, err = xorm.NewEngine(driverName, dsnName)
	if err != nil && err.Error() != "" {
		log.Fatal(err)
	}
	DbEngine.ShowSQL(false)
	//设置数据库连接数
	DbEngine.SetMaxOpenConns(10)
	//自动创建数据库
	DbEngine.Sync(new(User), new(Contact))

	fmt.Println("init database ok!")
}
