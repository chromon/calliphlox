package main

import (
	"calliphlox/dialect"
	"calliphlox/schema"
	"testing"
)

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T) {
	s := schema.Parse(&User{}, TestDial)
	if s.Name != "User" || len(s.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if s.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}