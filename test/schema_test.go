package main

import (
	"calliphlox/dialect"
	"calliphlox/schema"
	"testing"
)

type User struct {
	Id int `orm:"PRIMARY KEY"`
	Name string
	Age  int
}

var TestDial, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T) {
	s := schema.Parse(&User{}, TestDial)
	if s.Name != "User" || len(s.Fields) != 3 {
		t.Fatal("failed to parse User struct")
	}
	if s.GetField("Id").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}