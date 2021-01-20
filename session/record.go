package session

import (
	"calliphlox/clause"
	"errors"
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
	// 测试 hook
	s.CallMethod(BeforeQuery, nil)

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

		// 测试 hook
		s.CallMethod(AfterQuery, dest.Addr().Interface())

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

// update 接受 2 种入参，平铺开来的键值对和 map 类型的键值对
// s.Update(map[string]interface{}{"Name": "a", "Age": 18})
// s.Update("Name", "a", "Age", 18, ...)
func (s *Session) Update(kv ...interface{}) (int64, error) {
	// 因为 generator 接受的参数是 map 类型的键值对
	// 因此 Update 方法会动态地判断传入参数的类型
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		// kv 不是 map 类型，则会自动转换
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i + 1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// delete
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// count
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// Limit
// 通过链式(chain)操作，支持查询条件(where, order by, limit 等)的叠加
// 链式调用的原理也，某个对象调用某个方法后，将该对象的引用/指针返回，即可以继续调用该对象的其他方法
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// SQL 语句只返回一条记录
// u := &User{}
// _ = s.OrderBy("Age DESC").First(u)
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	// 获取切片的单个元素的类型
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}