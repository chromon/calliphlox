package main

import (
	"calliphlox"
	"fmt"

)

func main() {
	engine, _ := orm.NewEngine("mysql", "root:root@/calliphlox?charset=utf8")
	defer engine.Close()

	s := engine.NewSession()

	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(name varchar(20));").Exec()
	_, _ = s.Raw("CREATE TABLE User(name varchar(20));").Exec()
	result, _ := s.Raw("INSERT INTO User(`name`) values (?), (?)", "Ellery", "Feium").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}