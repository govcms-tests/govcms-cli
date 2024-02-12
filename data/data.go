package data

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DB_PATH = "/Users/jackfuller/dev/build/test/govcms.db"

var db *sql.DB

func Connect() {

	err := OpenDatabase()
	if err != nil {
		log.Fatal("Unable to connect to DB")
		return
	}
	SyncInstallations()
}

func OpenDatabase() error {
	var err error

	db, err = sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return err
	}

	return db.Ping()
}

func SyncInstallations() {
	listOfPaths := GetListOfPaths()
	for _, path := range listOfPaths {
		RemovePathIfMissing(path)
	}
}

func CreateTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS installations (
    	"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    	"name" TEXT UNIQUE NOT NULL,
    	"path" TEXT NOT NULL,
    	"type" TEXT NOT NULL
	);`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec()
	log.Println("Installations table created")
}

func InsertInstall(name string, path string, installType string) {
	insertInstallSQL := `INSERT INTO installations(name, path, type) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertInstallSQL)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec(name, path, installType)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Inserted installation successfully!")
}

func RemoveInstall(path string) {
	deleteSQL := `DELETE FROM installations where path=?`
	statement, err := db.Prepare(deleteSQL)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec(path)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetInstallPath(name string) string {
	var path string

	selectInstallSQL := `SELECT path FROM installations WHERE name=?`
	statement, err := db.Prepare(selectInstallSQL)
	if err != nil {
		log.Fatalln(err)
	}
	err = statement.QueryRow(name).Scan(&path)
	if err != nil {
		log.Fatalln(err)
	}

	return path
}

func GetListOfPaths() []string {
	var listOfPaths []string

	selectSQL := `SELECT path FROM installations`
	statement, err := db.Prepare(selectSQL)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return listOfPaths
	}

	rows, err := statement.Query()

	defer rows.Close()

	for rows.Next() {
		var path string
		rows.Scan(&path)
		listOfPaths = append(listOfPaths, path)
	}
	return listOfPaths
}

func RemovePathIfMissing(path string) {
	if exists, _ := DirExists(path); !exists {
		RemoveInstall(path)
	}
}

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
