package main

import (
	"calliphlox"
	"testing"
)

func TestSession_CreateTable(t *testing.T) {
	engine, _ := orm.NewEngine("mysql", "root:root@/calliphlox?charset=utf8")
	defer engine.Close()

	s := engine.NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.TableExist() {
		t.Fatal("Failed to create table User")
	}
}