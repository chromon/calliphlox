# Calliphlox

Calliphlox is a simple ORM for Go.

## Features

* Object <=> Table Mapping Support

* Chainable APIs

* Transaction Support

* MySQL, Sqlite3 dialect Support


## Quick Start

* Create Engine

```Go
engine, _ := orm.NewEngine(driverName, dataSourceName)
defer engine.Close()
```

* Define a struct

```Go
type User struct {
	Id int `orm:"PRIMARY KEY"`
	Name string
	Age  int
}
```

* Get Session Schema parse Object to db table

```GO
s := engine.NewSession().Model(&User{})
```

* Create table
```GO
_ := s.DropTable()
_ := s.CreateTable()
```

* Insert data
```GO
affected, err := s.Insert(user)
if err != nil || affected != 1 {
     t.Fatal("failed to create record")
}
```

 ## LICENSE
 Calliphlox is distributed under the terms of the GPL-3.0 License.
