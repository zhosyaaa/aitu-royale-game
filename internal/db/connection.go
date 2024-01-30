package db

import (
	"auth/internal/config"
	logger "auth/pkg/logger"
	"database/sql"
	"fmt"
	"sync"

	// Import the PostgreSQL driver
	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func initializeDB(database config.Database) error {
	var errInit error
	dbOnce.Do(func() {
		// Use the driver name "postgres" instead of "postgresql"
		dbConnString := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			database.User,
			database.Password,
			database.Host,
			database.Port,
			database.Name,
			database.Sslmode,
		)
		// Open the connection with the correct driver name
		db, errInit = sql.Open("postgres", dbConnString)
		if errInit != nil {
			logger.GetLogger().Fatal("Error connecting to the database:", errInit)
			return
		}
		errPing := db.Ping()
		if errPing != nil {
			logger.GetLogger().Fatal("Error pinging the database:", errPing)
			return
		}
		logger.GetLogger().Info("Connected to the database")
	})
	return errInit
}

func GetDBInstance(database config.Database) (*sql.DB, error) {
	db = nil
	var errGetDB error
	if db == nil {
		if err := initializeDB(database); err != nil {
			errGetDB = err
		}
	}
	return db, errGetDB
}
