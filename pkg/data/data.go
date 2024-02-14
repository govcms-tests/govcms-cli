package data

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/govcms-tests/govcms-cli/pkg/settings"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var db *sql.DB

func Initialise() {
	Connect()
	CreateTables()
	SyncInstallations()
	fmt.Println("Using database " + getDatabasePath())
}

func Connect() error {
	err := OpenDatabase()
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}
	return err
}

func OpenDatabase() error {
	createDatabaseIfNotExist()
	var err error
	db, err = sql.Open("sqlite3", getDatabasePath())
	if err != nil {
		log.Fatal(err)
	}
	return db.Ping()
}

func createDatabaseIfNotExist() {
	db, _ := sql.Open("sqlite3", getDatabasePath())
	defer db.Close()

	dbName := ".govcms.db"
	db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
}

func getDatabasePath() string {
	config, _ := settings.LoadConfig()
	return config.Database
}

func CreateTables() {
	CreateInstallationTables()
}

func CreateInstallationTables() {
	if TableExists("installations") {
		return
	}
	log.Println("Creating table 'installations'")
	createTableSQL := `CREATE TABLE IF NOT EXISTS installations (
    	"name" TEXT NOT NULL,
    	"path" TEXT PRIMARY KEY NOT NULL,
    	"type" INTEGER NOT NULL
	);`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec()
	log.Println("Installations table created")
}

func SyncInstallations() {
	listOfPaths := GetListOfPaths()
	for _, path := range listOfPaths {
		RemovePathIfMissing(path)
	}
}

func InsertInstallation(install Installation) {
	if InstallationExists(install) {
		return
	}
	insertInstallSQL := `INSERT INTO installations(name, path, type) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertInstallSQL)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec(install.Name, install.Path, install.Resource)
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

func InsertInstallations(installs []Installation) {
	for _, install := range installs {
		InsertInstallation(install)
	}
}

func InstallationExists(install Installation) bool {
	query := `SELECT path FROM installations WHERE path = ?`
	err := db.QueryRow(query, install.Path).Scan(&install.Path)
	if err == nil {
		return true
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Real error happened
		log.Print(err)
	}
	// No row was found
	return false
}

func TableExists(table string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
	err := db.QueryRow(query, table).Scan(&table)
	if err == nil {
		return true
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Real error happened
		log.Print(err)
	}
	// No row was found
	return false
}
