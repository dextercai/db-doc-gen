package doc

import (
	"fmt"
	"github.com/dextercai/db-doc-gen/model"
)

func genTableCol(cols []model.Column) []string {
	var tableMd []string
	//tableMd = append(tableMd, "## 列定义")
	tableMd = append(tableMd, "| 列名 | 类型 | KEY | 可否为空 | 默认值 | 注释 | 备注 |")
	tableMd = append(tableMd, "| ---- | ---- | ---- | ---- | ---- | ---- | ---- |")

	for j := range cols {
		tableMd = append(tableMd, fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |",
			cols[j].ColName, cols[j].ColType, cols[j].ColKey, cols[j].IsNullable, cols[j].ColDefault, cols[j].ColComment, ""))
	}
	return tableMd
}

func genTableIdx(idxs []model.Index) []string {
	var tableIdxMd []string
	//tableIdxMd = append(tableIdxMd, "## 索引信息")
	tableIdxMd = append(tableIdxMd, "| 索引名 | 关联字段 | 联合索引顺序 | 索引类型 | 可重复（NOT UNIQUE） | 注释 | 备注 |")
	tableIdxMd = append(tableIdxMd, "| ---- | ---- | ---- | ---- | ---- | ---- | ---- |")

	for j := range idxs {
		tableIdxMd = append(tableIdxMd, fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |",
			idxs[j].IndexName, idxs[j].ColName, idxs[j].SeqIndex, idxs[j].IndexType, idxs[j].IsNotUnique, idxs[j].Comment, ""))
	}
	return tableIdxMd
}

func genTableSqlArea() []string {
	var tableSQLMd []string
	//tableSQLMd = append(tableSQLMd, "## 数据库定义SQL")
	tableSQLMd = append(tableSQLMd, "```")
	tableSQLMd = append(tableSQLMd, "// TODO：待填充数据库定义SQL")
	tableSQLMd = append(tableSQLMd, "```")
	return tableSQLMd
}
