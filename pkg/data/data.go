package data

import (
	"database/sql"
	"github.com/govcms-tests/govcms-cli/pkg/settings"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// TODO Use an ORM such as SQLBoiler or GORP

//	type LocalStorage struct {
//		db *sql.DB
//	}
func CreateDatabaseIfNotExist() (*sql.DB, error) {
	dbPath := getDatabasePath()
	createFileIfNotExist(dbPath)
	db, err := sql.Open("sqlite3", dbPath)

	checkError(err)
	return db, db.Ping()
}

func getDatabasePath() string {
	config, _ := settings.LoadConfig()
	return config.Database
}

func createFileIfNotExist(path string) {
	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		checkError(err)
		file.Close()
	}
}

//	func Initialise(db *sql.DB) LocalStorage {
//		local, _ := Connect(db)
//		local.CreateTables()
//		local.SyncInstallations()
//		return local
//	}
//
//	func Connect(db *sql.DB) (LocalStorage, error) {
//		local := LocalStorage{db: db}
//		return local, db.Ping()
//	}
//
//	func (local *LocalStorage) CreateTables() {
//		local.CreateInstallationTables()
//	}
//
//	func (local *LocalStorage) CreateInstallationTables() {
//		if local.TableExists("installations") {
//			return
//		}
//		log.Println("Creating table 'installations'")
//		createTableSQL := `CREATE TABLE IF NOT EXISTS installations (
//	   	"name" TEXT UNIQUE NOT NULL,
//	   	"path" TEXT PRIMARY KEY NOT NULL,
//	   	"type" INTEGER NOT NULL
//		);`
//
//		statement, err := local.db.Prepare(createTableSQL)
//		checkError(err)
//		statement.Exec()
//		log.Println("Installations table created")
//	}
//
//	func (local *LocalStorage) SyncInstallations() {
//		listOfPaths := local.GetListOfPaths()
//		for _, path := range listOfPaths {
//			local.RemovePathIfMissing(path)
//		}
//	}
//
//	func (local *LocalStorage) InsertInstallation(install Installation) error {
//		if local.DoesInstallationExist(install) {
//			return errors.New("an installation already exists with this name")
//		}
//		insertInstallSQL := `INSERT INTO installations(name, path, type) VALUES (?, ?, ?)`
//		statement, err := local.db.Prepare(insertInstallSQL)
//		checkError(err)
//		_, err = statement.Exec(install.Name, install.Path, install.Resource)
//		checkError(err)
//		log.Println("Inserted installation successfully!")
//		return nil
//	}
//
//	func (local *LocalStorage) RemoveInstallFromPath(path string) {
//		deleteSQL := `DELETE FROM installations where path=?`
//		statement, err := local.db.Prepare(deleteSQL)
//		checkError(err)
//		_, err = statement.Exec(path)
//		checkError(err)
//	}
//
//	func (local *LocalStorage) RemoveInstallFromName(path string) error {
//		deleteSQL := `DELETE FROM installations where name=?`
//		statement, err := local.db.Prepare(deleteSQL)
//		checkError(err)
//		_, err = statement.Exec(path)
//		if errors.Is(err, sql.ErrNoRows) {
//			return fmt.Errorf("no installation found with that name")
//		}
//		return err
//	}
//
//	func (local *LocalStorage) GetInstallPath(name string) (string, error) {
//		var path string
//
//		selectInstallSQL := `SELECT path FROM installations WHERE name=?`
//		statement, err := local.db.Prepare(selectInstallSQL)
//		checkError(err)
//		err = statement.QueryRow(name).Scan(&path)
//		if err != nil {
//			return "", err
//		}
//		return path, nil
//	}
//
//	func (local *LocalStorage) GetListOfPaths() []string {
//		var listOfPaths []string
//
//		selectSQL := `SELECT path FROM installations`
//		statement, err := local.db.Prepare(selectSQL)
//		if err != nil {
//			fmt.Println("Error executing query:", err)
//			return listOfPaths
//		}
//
//		rows, err := statement.Query()
//
//		defer rows.Close()
//
//		for rows.Next() {
//			var path string
//			rows.Scan(&path)
//			listOfPaths = append(listOfPaths, path)
//		}
//		return listOfPaths
//	}
//
//	func (local *LocalStorage) RemovePathIfMissing(path string) {
//		if exists, _ := DirExists(path); !exists {
//			local.RemoveInstallFromPath(path)
//		}
//	}
//
//	func DirExists(path string) (bool, error) {
//		_, err := os.Stat(path)
//		if err == nil {
//			return true, nil
//		}
//		if errors.Is(err, os.ErrNotExist) {
//			return false, nil
//		}
//		return false, err
//	}
//
//	func (local *LocalStorage) InsertInstallations(installs []Installation) {
//		for _, install := range installs {
//			local.InsertInstallation(install)
//		}
//	}
//
//	func (local *LocalStorage) DoesInstallationExist(install Installation) bool {
//		query := `SELECT path FROM installations WHERE path = ?`
//		err := local.db.QueryRow(query, install.Path).Scan(&install.Path)
//		if err == nil {
//			return true
//		}
//		if !errors.Is(err, sql.ErrNoRows) {
//			// Real error happened
//			log.Print(err)
//		}
//		// No row was found
//		return false
//	}
//
//	func (local *LocalStorage) TableExists(table string) bool {
//		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
//		err := local.db.QueryRow(query, table).Scan(&table)
//		if err == nil {
//			return true
//		}
//		if !errors.Is(err, sql.ErrNoRows) {
//			// Real error happened
//			log.Print(err)
//		}
//		// No row was found
//		return false
//	}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
