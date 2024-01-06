package db

import (
	"auth/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func initializeDB(database config.Database) error {
	var err error
	dbOnce.Do(func() {
		dbConnString := fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			database.User,
			database.Password,
			database.Host,
			database.Port,
			database.Name,
			database.Sslmode,
		)
		db, err = gorm.Open(postgres.Open(dbConnString), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			log.Println("Error connecting to the database:", err)
			return
		}
		log.Println("Connected to the database")
	})
	return err
}

func GetDBInstance(database config.Database) (*gorm.DB, error) {
	db = nil
	if db == nil {
		if err := initializeDB(database); err != nil {
			return nil, err
		}
	}
	return db, nil
}
