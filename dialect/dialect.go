package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

// 方言, 用于处理数据库间不同差异，实现最大程度的复用和解耦
type Dialect interface {
	// 将 Go 语言的类型利用反射转换为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string

	// 返回某个表是否存在的 sql 语句
	TableExistSQL(tableName string) (string, []interface{})
}

// 注册方言实例
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 获取方言实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}