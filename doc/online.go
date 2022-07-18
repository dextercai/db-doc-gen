package doc

import (
	"fmt"
	"github.com/dextercai/db-doc-gen/model"
	"github.com/dextercai/db-doc-gen/util"
	"log"
	"net/http"
	"path"
	"strings"
)

func CreateOnlineDoc(docPath string, dbInfo model.DbInfo, tables []model.Table, listen string) {
	var sidebar []string
	var readme []string
	readme = append(readme, fmt.Sprintf("# %s 数据库文档", dbInfo.DbName))
	readme = append(readme, "### 基础信息")
	readme = append(readme, "| 数据库名称 | 版本 | 字符集 | 排序规则 |")
	readme = append(readme, "| ---- | ---- | ---- | ---- |")
	readme = append(readme, fmt.Sprintf("| %s | %s | %s | %s |", dbInfo.DbName, dbInfo.Version, dbInfo.Charset, dbInfo.Collation))
	readme = append(readme, "")
	readme = append(readme, "*created by github.com/dextercai/db-doc-gen*")

	for i := range tables {
		sidebar = append(sidebar, fmt.Sprintf("* [%s(%s)](%s.md)", tables[i].TableName, tables[i].TableComment, tables[i].TableName))

		var tableMd []string
		tableMd = append(tableMd, "## 字段定义")
		tableMd = append(tableMd, genTableCol(tables[i].ColList)...)
		tableStr := strings.Join(tableMd, "\r\n")

		var tableIdx []string
		tableIdx = append(tableIdx, "## 索引定义")
		tableIdx = append(tableIdx, genTableIdx(tables[i].IdxList)...)
		tableIdxStr := strings.Join(tableIdx, "\r\n")

		var tableSql []string
		tableSql = append(tableSql, "## DDL")
		tableSql = append(tableSql, genTableSqlArea(tables[i].Ddl)...)
		tableSqlStr := strings.Join(tableSql, "\r\n")

		pageContent := fmt.Sprintf("# %s(%s)\r\n", tables[i].TableName, tables[i].TableComment) +
			tableStr + "\r\n\r\n" + tableIdxStr + "\r\n\r\n" + tableSqlStr

		util.WriteToFile(path.Join(docPath, tables[i].TableName+".md"), pageContent)
	}

	readmeStr := strings.Join(readme, "\r\n")
	util.WriteToFile(path.Join(docPath, "README.md"), readmeStr)
	sidebarStr := strings.Join(sidebar, "\r\n")
	util.WriteToFile(path.Join(docPath, "_sidebar.md"), sidebarStr)
	util.WriteToFile(path.Join(docPath, "index.html"), docsifyHTML)
	util.WriteToFile(path.Join(docPath, ".nojekyll"), "")
	fmt.Println("文档生成成功！")
	runServer(docPath, listen)
}

// runServer run http static server
func runServer(dir string, listen string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	fmt.Println("http服务监听地址: " + listen)
	fmt.Println("文档存放地址: " + http.Dir(dir))

	log.Fatal(http.ListenAndServe(listen, nil))
}
