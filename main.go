package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var server = "localhost"
var port = 1433
var user = "admin"
var password = "admin"
var database = "employees"

func main() {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	db, err := gorm.Open("mssql", connString)
	if err != nil {
		fmt.Println(err)
	}

	// 防止 gorm 產生 SQL 時，自動在 table 加上 s
	// https://gorm.io/docs/conventions.html#Pluralized-Table-Name
	db.SingularTable(true)

	var id int
	var name string

	db.Debug().Table("check_in_record").Select("Id, Name").Row().Scan(&id, &name)

	fmt.Println(id, name)
}

/*****下面為測試的code****/
// package main

// import (
// 	"fmt"
// 	"strings"
// )

// func WordCount(s string) map[string]int{

// 	var a[] string = strings.Fields(s)
// 	var m map[string]int
// 	m = make(map[string]int)

// 	fmt.Println(a)
// 	for i := 0; i<len(a);i++ {
// 		//fmt.Println(a[i])
// 		word := a[i]
// 		m[word] = len(a[i])
// 		fmt.Println(m)
// 	}
// 	return m
// }

// func main() {

// 	s := "foo   bar    baz"
// 	fmt.Println(s)

// 	var a[] string = strings.Fields(s)
// 	fmt.Println(cap(a),len(a))

// 	var m map[string]int
// 	m = WordCount("a b c dddddd ")

// 	fmt.Println(m)
// 	//fmt.Println(WordCount("a b c d "))
// 	// fmt.Printf("%q",strings.Fields(" foo   bar    baz"))

// }
