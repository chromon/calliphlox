package session

import (
	"calliphlox/clause"
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