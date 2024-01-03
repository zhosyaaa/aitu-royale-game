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
	dbConf config.Database
)

func InitializeDB(config config.Database) error {
	dbConf = config
	return initializeDB()
}

func initializeDB() error {
	var err error
	dbOnce.Do(func() {
		dbConnString := fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			dbConf.User,
			dbConf.Password,
			dbConf.Host,
			dbConf.Port,
			dbConf.Name,
			dbConf.Sslmode,
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

func GetDBInstance() (*gorm.DB, error) {
	if db == nil {
		if err := initializeDB(); err != nil {
			return nil, err
		}
	}
	return db, nil
}
