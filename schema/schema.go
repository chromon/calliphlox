package schema

import (
	"calliphlox/dialect"
	"go/ast"
	"reflect"
)

// 对象 (object) 和表 (table) 的转换

// 数据库中的列属性
type Field struct {
	// 数据表字段名称
	Name string

	// 数据表字段类型，是数据库实际类型
	Type string

	// 数据表约束条件 (非空、主键等通过 go 语言 Tag 实现)
	Tag string
}

// 数据库中的表属性
type Schema struct {
	// 数据库表对应的实体类对象
	Model interface{}

	// 数据库表名
	Name string

	// 数据库表字段属性数组
	Fields []*Field

	// 数据库表字段名称数组
	FieldNames []string

	// 数据库表字段名称 (FieldName) 与字段属性 (Field) 的映射关系
	fieldMap map[string]*Field
}

// 由数据库表字段名称查询字段属性对象
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 将对象解析为 Schema 实例
func Parse(dest interface{}, dialect dialect.Dialect) *Schema {
	// 通过反射获取对象类型
	// ValueOf() 返回传入参数的值
	// TypeOf() 返回传入参数的类型
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	// 构造数据库表 Schema 对象
	schema := &Schema{
		Model: dest,
		Name: modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	// NumField() 返回实例字段的个数
	for i := 0; i < modelType.NumField(); i++ {
		// 获取第 i 个成员变量
		p := modelType.Field(i)

		// func IsExported(name string) bool 否为导出的 Go 符号（即是否以大写字母开头）
		// Anonymous bool 是否匿名字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			// 构造数据库表字段
			field := &Field {
				// 字段名
				Name: p.Name,
				// 字段类型，当前为数据库实际类型
				Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			// 设置 tag 获取约束条件
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}

			// 字段属性实例
			schema.Fields = append(schema.Fields, field)
			// 字段名
			schema.FieldNames = append(schema.FieldNames, p.Name)
			// 字段名与字段属性实例映射关系
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}