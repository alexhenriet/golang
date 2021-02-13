// http://go-database-sql.org/
// https://github.com/go-sql-driver/mysql
// https://github.com/go-sql-driver/mysql/wiki/Examples
// https://golang.org/pkg/database/sql/

package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"regexp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	id int
	firstname string
	name string
}

func main() {
	argv := os.Args
	if len(argv) != 2 {
		fmt.Printf("Error: use %s .env_file\n", argv[0])
		os.Exit(1)
	}
	_ = godotenv.Load(argv[1])
	databaseURL := fixDsn(os.Getenv("DATABASE_URL"))

	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()
	query := "SELECT 1 as id, \"john\" as firstname, \"doe\" as name FROM dual;"
	results, err := db.Query(query)
    if err != nil {
        panic(err.Error())
    }
    for results.Next() {
        var user User
        err = results.Scan(&user.id, &user.firstname, &user.name)
        if err != nil {
             panic(err.Error())
        }
        fmt.Printf("%v\n", user)
    }
}

// Convert DSN FROM Symfony dotenv syntax for Doctrine TO go-database-sql syntax
func fixDsn(original string) string {
	re := regexp.MustCompile(`mysql://(.+):(.+)@(.+):(.+)/(.+)`)
	matches := re.FindAllStringSubmatch(original, -1)
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", matches[0][1], matches[0][2], matches[0][3], matches[0][4], matches[0][5])
}