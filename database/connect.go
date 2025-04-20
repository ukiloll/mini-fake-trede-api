package database

import (
	"database/sql"
	"fmt"
	"github/ukilolll/trade/pkg"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var(
	_ = pkg.LoadEnv()
)

var(
	DB_USERNAME=os.Getenv("DB_USERNAME")
	DB_PASSWORD=os.Getenv("DB_PASSWORD")
	DB_HOST=os.Getenv("DB_HOST")
	DB_NAME=os.Getenv("DB_NAME")

)
func Connect() *sql.DB{
	var dsn = fmt.Sprintf("%v:%v@tcp(%v)/%v",DB_USERNAME,DB_PASSWORD,DB_HOST,DB_NAME)

	conn,err := sql.Open("mysql",dsn)
	if err != nil {
		log.Panic(err)
	}
	return conn
}