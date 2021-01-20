package clause

import (
	"fmt"
	"strings"
)

// 生成子句
// 子句生成方法类型
type generator func(values ...interface{}) (string, []interface{})

// 子句类型与相关方法对应 map
var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)

	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// 根据变量个数生成同等数量占位符
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// 插入子句
// INSERT INTO tableName (...)
func _insert(values ...interface{}) (string, []interface{}) {
	// 表名
	tableName := values[0]
	// 待插入字段名
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName,fields), []interface{}{}
}

// values 子句，可以同时插入多条
// VALUES (?, ? ...), (?, ? ...) ... (?, ? ...)
func _values(values ...interface{}) (string, []interface{}) {
	// 插入字段占位符
	var bindStr string
	// sql 语句
	var sql strings.Builder
	// 实际的待插入数据切片，与占位符对应
	var vars []interface{}

	sql.WriteString("VALUES ")
	// 遍历插入数据，可以同时插入多条
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i + 1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// select 子句
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// limit 子句
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// where 子句
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

// orderBy 子句
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// update 子句
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	// map 类型，表示待更新的键值对
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k + " = ?")
		vars = append(vars, v)
	}

	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

// delete 子句
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// count 子句
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}