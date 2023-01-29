# db-doc-gen
我习惯性的在项目中直接通过IDL或者类似的手段直接对数据库表进行操作；这也是敏捷开发时常遇到的。

但避免不了在项目沟通时依旧需要对数据库字段进行说明。

故简单编写了一个适用于MySQL/MariaDB的文档生成。

## TODO 
- [x] MySQL/MariaDB
- [x] 结构导出为SQL

## HOW TO INSTALL
```shell
go install github.com/dextercai/db-doc-gen
```

## HOW TO USE
```
github.com/dextercai/db-doc-gen (unknown)

Usage of db-doc-gen:
  -database string
        数据库名 (default "test")
  -db-type string
        数据库类型：mysql (default "mysql")
  -doc-serve string
        在线文档服务地址 (default ":8080")
  -doc-type string
        文档生成类型：online、offline (default "online")
  -host string
        数据库地址 (default "127.0.0.1")
  -password string
        数据库密码 (default "123456")
  -port int
        数据库端口 (default 3306)
  -username string
        数据库用户名 (default "admin")
```
