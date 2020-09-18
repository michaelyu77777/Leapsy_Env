package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var server = "localhost"
var port = 1433
var user = "admin"
var password = "admin"
var database = "employees"

func main() {

	// Connect to database
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	fmt.Printf("Connected!\n")
	defer conn.Close()

	// Create employee
	createID, err := CreateCheckInRecord(conn, "Jake", "2020-09-18", "t", "t", "2020/08/30", "軟體", "軟體工程師")
	if err != nil {
		log.Fatal("CreateEmployee failed:", err.Error())
	}
	fmt.Printf("成功建立= %d.\n", createID)

	// connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
	// 	server, user, password, port, database)

	// db, err := gorm.Open("mssql", connString)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // 防止 gorm 產生 SQL 時，自動在 table 加上 s
	// // https://gorm.io/docs/conventions.html#Pluralized-Table-Name
	// db.SingularTable(true)

	// var id int
	// var name string

	// db.Debug().Table("check_in_record").Select("Id, Name").Row().Scan(&id, &name)

	// fmt.Println(id, name)

	// // Create employee
	// createID, err := Create_CheckInRecord(conn, "Jake", "United States")
	// if err != nil {
	// 	log.Fatal("CreateEmployee failed:", err.Error())
	// }
	// fmt.Printf("Inserted ID: %d successfully.\n", createID)
}

// CreateCheckInRecord return (int64, error)
func CreateCheckInRecord(db *sql.DB, name string, checkInTime string, pic string, leaveType string, date string, department string, position string) (int64, error) {

	// tsql := fmt.Sprintf("INSERT INTO check_in_record(name,check_in_time,pic,leave_type,date,department,position) VALUES ('%s','%s','%s','%s','%s','%s','%s');",
	// 	name,
	// 	checkInTime,
	// 	//pic,
	// 	nil,
	// 	leaveType,
	// 	date,
	// 	department,
	// 	position)
	// _, err := db.Exec(tsql)

	var newPic []string
	fmt.Println(newPic)

	_, err := db.Exec("INSERT INTO check_in_record(name,check_in_time,pic,leave_type,date,department,position) VALUES (?,?,?,?,?,?,?)", name, checkInTime, NewNullString(pic), NewNullString(leaveType), date, department, position)

	if err != nil {
		fmt.Println("Error inserting new row: " + err.Error())
		return -1, err
	}

	return 1, err
}

// NewNullString returns sql.NullString
// 工具:若string為空字串，可以轉成nil (若傳入為空:回傳nil, 若非空:回傳一個物件 sql.NullString 但在取的時候 很神奇 會自動取到 String這個參數的值)
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
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
