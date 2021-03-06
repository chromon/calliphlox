package main

import (
	"calliphlox/clause"
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var cla clause.Clause
	cla.Set(clause.LIMIT, 3)
	cla.Set(clause.SELECT, "User", []string{"*"})
	cla.Set(clause.WHERE, "Name = ?", "a")
	cla.Set(clause.ORDERBY, "Age ASC")
	sql, vars := cla.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"a", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}