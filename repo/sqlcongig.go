package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // ✅ this registers the MySQL driver
)

func Connectdb(dbname string) (*sql.DB, error) {
	coneectionString := "root:1@tcp(127.0.0.1:3306)/" + dbname
	db, err := sql.Open("mysql", coneectionString)
	if err != nil {
		panic(err)
	}
	fmt.Println("sql coneection ")
	return db, nil

}
