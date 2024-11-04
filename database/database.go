package database

import (
	"fmt"
	"lock/config"
	"lock/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// func InitDB(cfg config.Config) {

// 	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s",
// 		cfg.DBHost,
// 		cfg.DBUser,
// 		cfg.DBName,
// 		cfg.DBPort,
// 		cfg.DBPassword,
// 	)
// 	var err error

// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

// 	if err != nil {

// 		log.Fatal("Failid to connec database ", err)
// 	}

// 	sqlDB, err := DB.DB()

// 	if err != nil {

// 		log.Fatal("Failed to get database object ", err)
// 	}

// 	err = sqlDB.Ping()

// 	if err != nil {

// 		log.Fatal("Falied  to ping the database ")
// 	}

// 	log.Println("Databse connection was success fully ")
// }

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)

	db, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if dberr != nil {

		return nil, fmt.Errorf("faild to connect to database:%w", dberr)
	}

	DB = db
	DB.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Admin{})

	return DB, nil
}
