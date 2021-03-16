package main

import ( //our Go packages for this project
	"database/sql"
	"log"
	_"github.com/denisenkom/go-mssqldb"
)

// where our program begins
// where our program begins
func main() {
	//connect to the database
	db, err := sql.Open("sqlserver", "sqlserver://[username]:[password]@[host]?database=[database]&connection+timeout=30")
if err != nil {
	log.Fatal("Open connection failed:",err.Error())
	}

defer db.Close()


}

