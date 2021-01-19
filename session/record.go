package session

import (
	"calliphlox/clause"
	"reflect"
)

// 实现记录增删查改

// 插入
// s.Insert(u1, u2, ...)
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		// 对象解析为 schema
		table := s.Model(value).RefTable()
		// 调用对应的 insert generator，生成该对应的 SQL 语句
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// 将对象属性转为数据库参数顺序列表
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// 设置 value 子句
	s.clause.Set(clause.VALUES, recordValues...)
	// 构建完整 sql 语句
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Find 功能
// s.Find(&u)
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	// 获取切片的单个元素的类型
	destType := destSlice.Type().Elem()
	// 使用 reflect.New() 方法创建一个 destType 的实例，作为 Model() 的入参，映射出表结构 RefTable()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// 根据表结构，使用 clause 构造出 SELECT 语句
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	// 查询到所有符合条件的记录 rows
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// 遍历每一行记录
	for rows.Next() {
		// 利用反射创建 destType 的实例 dest
		dest := reflect.New(destType).Elem()
		var values []interface{}
		// 将 dest 的所有字段平铺开，构造切片 values
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 将该行记录每一列的值依次赋值给 values 中的每一个字段
		if err := rows.Scan(values...); err != nil {
			return err
		}
		// 将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}