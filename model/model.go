package model

type DbConfig struct {
	DbType   *string // 1. mysql
	DocType  *string // 1. online 2. offline
	DocServe *string
	Host     *string
	Port     *int
	User     *string
	Password *string
	Database *string
}

type DbInfo struct {
	Version   string
	Charset   string
	Collation string
	DbName    string
}

type Column struct {
	ColName    string
	ColType    string
	ColKey     string
	IsNullable string
	ColComment string
	ColDefault string
}

type Index struct {
	IndexName   string
	ColName     string
	SeqIndex    string
	IndexType   string
	Comment     string
	IsNotUnique string
}

type Table struct {
	TableName    string
	TableComment string
	ColList      []Column
	IdxList      []Index
	Ddl          TableCreateSql
}

type TableCreateSql struct {
	TableName string
	SqlLine   string
}
