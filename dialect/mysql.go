package dialect

import (
	"fmt"
	"reflect"
	"time"
)

// mysql 方言
type mysql struct{}

// 检测某个类型是否实现了某个接口的所有方法
// 将 nil 转成 *mysql 类型，并赋值给 Dialect，如果失败则表明 mysql 并没有实现 Dialect 全部方法
var _ Dialect = (*mysql)(nil)

// 注册 mysql 方言实例
func init() {
	RegisterDialect("mysql", &mysql{})
}

// 将 Go 语言的类型利用反射转换为 mysql 数据库的数据类型
// reflect.ValueOf() 获取指针对应的反射值
// reflect.Indirect() 获取指针指向的对象的反射值
// (reflect.Type).Name() 返回类名(字符串)
// (reflect.Type).Field(i) 获取第 i 个成员变量
// Value.Kind() 返回 value 的类型
func (m *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int32:
		return "integer"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int16:
		return "smallint"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint, reflect.Uint32:
		return "integer unsigned"
	case reflect.Uint8:
		return "tinyint unsigned"
	case reflect.Uint16:
		return "smallint unsigned"
	case reflect.Uint64:
		return "bigint unsigned"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}

	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// 获取 mysql 方言实例
func (m *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{} {
		tableName,
	}
	return "SHOW TABLES LIKE ?", args
}