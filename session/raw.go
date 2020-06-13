package session

import (
	"calliphlox/log"
	"database/sql"
	"strings"
)

// 负责实现与数据库交互（原生）
type Session struct {
	// 数据库连接
	db *sql.DB

	// sql 语句
	sql strings.Builder

	// sql 语句占位符
	sqlVars []interface{}
}

// 由数据库连接创建 session 对象
func New(db *sql.DB) *Session {
	return &Session{
		db: db,
	}
}

// 重置 sql 语句
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// 获取数据库连接
func (s *Session) DB() *sql.DB {
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
func (s *Session) QueryRows(rows *sql.Rows, err error) {
	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}