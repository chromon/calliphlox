package session

import (
	"calliphlox/clause"
	"calliphlox/dialect"
	"calliphlox/log"
	"calliphlox/schema"
	"database/sql"
	"strings"
)

// 负责实现与数据库交互（原生）
type Session struct {
	// 数据库连接
	db *sql.DB

	// 事务
	tx *sql.Tx

	// sql 语句
	sql strings.Builder

	// sql 语句占位符
	sqlVars []interface{}

	// 方言
	dialect dialect.Dialect

	// 数据库表对象
	refTable *schema.Schema

	// 子句
	clause clause.Clause
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// 验证是否实现了 CommonDB 接口
var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// 由数据库连接创建 session 对象
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db: db,
		dialect: dialect,
	}
}

// 重置 sql 语句
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// 获取数据库连接
func (s *Session) DB() CommonDB {
	// 当 tx 不为空时，则使用 tx 执行 SQL 语句，否则使用 db 执行 SQL 语句
	if s.tx != nil {
		return s.tx
	}

	return s.db
}

// 拼接 sql 语句
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// 执行原生 sql 语句
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)

	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// 从数据库中获取一行记录
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// 从数据库中读取多行记录
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}