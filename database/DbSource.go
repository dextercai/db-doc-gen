package database

import (
	"database/sql"
	"fmt"
	"github.com/dextercai/db-doc-gen/model"
	_ "github.com/go-sql-driver/mysql"
)

type DbSource struct {
	DB     *sql.DB
	Config model.DbConfig
}

func (this *DbSource) GetDbInfo() model.DbInfo {

	var (
		info       model.DbInfo
		rows       *sql.Rows
		err        error
		key, value string
	)

	// 数据库版本
	rows, err = this.DB.Query("select @@version;")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
	}
	info.Version = value
	// 字符集
	rows, err = this.DB.Query("show variables like '%character_set_server%';")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&key, &value)
	}
	info.Charset = value
	// 排序规则
	rows, err = this.DB.Query("show variables like 'collation_server%';")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&key, &value)
	}
	info.Collation = value

	return info
}

func (this DbSource) GetTableInfo() []model.Table {
	tables := make([]model.Table, 0)
	rows, err := this.DB.Query(this.getTableSQL())
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var table model.Table
	for rows.Next() {
		table.TableComment = ""
		rows.Scan(&table.TableName, &table.TableComment)
		if len(table.TableComment) == 0 {
			table.TableComment = table.TableName
		}
		tables = append(tables, table)
	}
	for i := range tables {
		tables[i].ColList = this.getColumnInfo(tables[i].TableName)
		tables[i].IdxList = this.getIndexInfo(tables[i].TableName)
		tables[i].Ddl = this.getTableCreateSqlInfo(tables[i].TableName)
	}
	return tables
}

func (this DbSource) getColumnInfo(tableName string) []model.Column {
	columns := make([]model.Column, 0)
	rows, err := this.DB.Query(this.getColumnSQL(tableName))
	if err != nil {
		fmt.Println(err)
	}
	var column model.Column
	for rows.Next() {
		rows.Scan(&column.ColName, &column.ColType, &column.ColKey, &column.IsNullable, &column.ColComment, &column.ColDefault)
		columns = append(columns, column)
	}
	return columns
}

func (this DbSource) getTableSQL() string {
	var sqlLine string
	if *this.Config.DbType == "mysql" {
		sqlLine = fmt.Sprintf(`
			select table_name    as TableName, 
			       table_comment as TableComment
			from information_schema.tables 
			where table_schema = '%s'
		`, *this.Config.Database)
	}

	return sqlLine
}

func (this *DbSource) getColumnSQL(tableName string) string {
	var sqlLine string
	if *this.Config.DbType == "mysql" {
		sqlLine = fmt.Sprintf(`
			select column_name as ColName,
			column_type        as ColType,
			column_key         as ColKey,
			is_nullable        as IsNullable,
			column_comment     as ColComment,
			column_default     as ColDefault
			from information_schema.columns 
			where table_schema = '%s' and table_name = '%s'
		`, *this.Config.Database, tableName)
	}

	return sqlLine
}

func (this *DbSource) getIndexSQL(tableName string) string {
	var sqlLine string
	if *this.Config.DbType == "mysql" {
		sqlLine = fmt.Sprintf(`
			SELECT INDEX_NAME as IndexName, 
			COLUMN_NAME as ColName, 
			SEQ_IN_INDEX as SeqIndex, 
			INDEX_TYPE as IndexType, 
			INDEX_COMMENT as Comment, 
			NON_UNIQUE as isNotUnique 
			from information_schema.STATISTICS
			where table_schema = '%s' and table_name = '%s'
		`, *this.Config.Database, tableName)
	}
	return sqlLine
}

func (this *DbSource) getTableCreateSql(tableName string) string {
	var sqlLine string
	if *this.Config.DbType == "mysql" {
		sqlLine = fmt.Sprintf(`
			SHOW CREATE TABLE %s.%s
		`, *this.Config.Database, tableName)
	}
	return sqlLine
}

func (this *DbSource) getIndexInfo(tableName string) []model.Index {
	idxs := make([]model.Index, 0)
	rows, err := this.DB.Query(this.getIndexSQL(tableName))
	if err != nil {
		fmt.Println(err)
	}
	var idx model.Index
	for rows.Next() {
		rows.Scan(&idx.IndexName, &idx.ColName, &idx.SeqIndex, &idx.IndexType, &idx.Comment, &idx.IsNotUnique)
		idxs = append(idxs, idx)
	}
	return idxs
}

func (this *DbSource) getTableCreateSqlInfo(tableName string) model.TableCreateSql {
	createSql := model.TableCreateSql{}
	rows, err := this.DB.Query(this.getTableCreateSql(tableName))
	if err != nil {
		fmt.Println(err)
	}
	// 期望只有一行
	for rows.Next() {
		rows.Scan(&createSql.TableName, &createSql.SqlLine)
	}
	return createSql
}
