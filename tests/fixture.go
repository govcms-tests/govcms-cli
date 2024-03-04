package tests

import (
	database "database-sqlc"
	"database/sql"
	"github.com/govcms-tests/govcms-cli/cmd"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
)

var db *sql.DB

func setUp() *cobra.Command {
	if _, err := os.Stat("test.db"); err == nil {
		_ = os.Remove("test.db")
	}
	db, _ := database.NewDatabase("test.db")
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), db)
	return rootCmd
}
