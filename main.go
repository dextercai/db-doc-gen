package main

import (
	"flag"
	"fmt"
	"github.com/dextercai/db-doc-gen/database"
	"github.com/dextercai/db-doc-gen/define"
	"github.com/dextercai/db-doc-gen/doc"
	"github.com/dextercai/db-doc-gen/model"
	"github.com/dextercai/db-doc-gen/util"
	"os"
	"path"
)

var dbConfig model.DbConfig

func main() {

	fmt.Printf("github.com/dextercai/db-doc-gen (%s)\n", define.VERSION)

	dbConfig = model.DbConfig{
		DbType:   flag.String("db-type", "mysql", "数据库类型：mysql"),
		DocType:  flag.String("doc-type", "online", "文档生成类型：online、offline"),
		DocServe: flag.String("doc-serve", ":8080", "在线文档服务地址"),
		Host:     flag.String("host", "127.0.0.1", "数据库地址"),
		Port:     flag.Int("port", 3306, "数据库端口"),
		User:     flag.String("username", "admin", "数据库用户名"),
		Password: flag.String("password", "123456", "数据库密码"),
		Database: flag.String("database", "test", "数据库名"),
	}

	flag.Parse()

	// 数据库文档生成配置

	fmt.Printf("解析到配置信息： \n")
	fmt.Printf("%s\t%s\n", "数据库类型", *dbConfig.DbType)
	fmt.Printf("%s\t%s\n", "文档生成类型", *dbConfig.DocType)
	fmt.Printf("%s\t%s\n", "文档服务地址", *dbConfig.DocServe)
	fmt.Printf("%s\t%s\n", "数据库地址", *dbConfig.Host)
	fmt.Printf("%s\t%d\n", "数据库端口", *dbConfig.Port)
	fmt.Printf("%s\t%s\n", "数据库用户名", *dbConfig.User)
	fmt.Printf("%s\t%s\n", "数据库密码", *dbConfig.Password)
	fmt.Printf("%s\t%s\n", "数据库名", *dbConfig.Database)

	db, err := database.InitDb(&dbConfig)
	if err != nil {
		panic(err)
	}

	dbInfo := db.GetDbInfo()
	dbInfo.DbName = *dbConfig.Database
	tables := db.GetTableInfo()

	var docPath string
	dir, _ := os.Getwd()
	if *dbConfig.DocType == "online" {
		docPath = path.Join(dir, "dist", dbInfo.DbName, "www")
		util.CreateDir(docPath)
		doc.CreateOnlineDoc(docPath, dbInfo, tables, *dbConfig.DocServe)
	} else {
		docPath = path.Join(dir, "dist", dbInfo.DbName)
		util.CreateDir(docPath)
		doc.CreateOfflineDoc(docPath, dbInfo, tables)
	}
}
