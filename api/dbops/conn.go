package dbops


import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err error
)
// init
func init(){
	dbConn,err=sql.Open(
		"mysql",
		"root:{PASSWORD}@tcp(localhost:3306)/viedoServices")
	if err!=nil{
		panic(err.Error())
	}
}
