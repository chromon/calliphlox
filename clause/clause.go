package clause

import "strings"

// SELECT col1, col2, ...
// 	FROM table_name
// 	WHERE [ conditions ]
// 	GROUP BY col1
// 	HAVING [ conditions ]

// 拼接子句
type Clause struct {
	sql map[Type]string
	sqlVars map[Type][]interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// 根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句
// Set(clause.SELECT, "User", []string{"name", "age"})
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}

	// return (string, []interface{})
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// 根据传入的 Type 的顺序，构造出最终的 SQL 语句
// sql, vars := cla.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}

	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}