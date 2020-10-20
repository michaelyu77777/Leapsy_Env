package main

import (
	_ "github.com/denisenkom/go-mssqldb"

	"./LeapsyPackages/servers"
)

// main - 主程式
func main() {
	servers.StartECServer() // 啟動環控伺服器
}
