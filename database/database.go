package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var Database *sql.DB

func Connect() {
	var err error

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)
	Database, err = sql.Open("postgres", dsn)

	if err != nil {
		panic(err)
	}
}

func CreateTables() error {
	basepath := "./"
	cwd, _ := os.Getwd()

	dirs := strings.Split(cwd, "/")
	lastIndex := len(dirs) - 1

	if dirs[lastIndex] != "be" {
		basepath = "../"
	}

	if err := CreateTableUsers(basepath); err != nil {
		return err
	}

	if err := CreateTableEntries(basepath); err != nil {
		return err
	}

	return nil
}

func CreateTableUsers(basepath string) error {
	sqlFile, err := os.Open(basepath + "database/sql/create_table_users.sql")
	if err != nil {
		return err
	}
	defer sqlFile.Close()

	stat, err := sqlFile.Stat()
	if err != nil {
		return err
	}

	buf := make([]byte, stat.Size())

	if _, err := sqlFile.Read(buf); err != nil {
		return err
	}

	if _, err := Database.Exec(string(buf)); err != nil {
		return err
	}

	return nil
}

func CreateTableEntries(basepath string) error {
	sqlFile, err := os.Open(basepath + "database/sql/create_table_entries.sql")
	if err != nil {
		return err
	}
	defer sqlFile.Close()

	stat, err := sqlFile.Stat()
	if err != nil {
		return err
	}

	buf := make([]byte, stat.Size())

	if _, err := sqlFile.Read(buf); err != nil {
		return nil
	}

	if _, err := Database.Exec(string(buf)); err != nil {
		return err
	}

	return nil
}
