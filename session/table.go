package session

import (
	"calliphlox/log"
	"calliphlox/schema"
	"fmt"
	"reflect"
	"strings"
)

// Session 中数据库表操作相关

// 赋值 refTable
func (s *Session) Model(value interface{}) *Session {

	// 如果传入的结构体名称不发生变化，则不会更新 refTable 的值
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		// 解析实体对象
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// 返回 refTable 对象
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// 创建数据库表
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string

	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	desc := strings.Join(columns, ",")
	fmt.Println(desc)
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()

	return err
}

// 删除数据库表
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.refTable.Name)).Exec()
	return err
}

// 判断数据库表是否存在
func (s *Session) TableExist() bool {
	sql, values := s.dialect.TableExistSQL(s.refTable.Name)
	row := s.Raw(sql, values...).QueryRow()

	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.refTable.Name
}

