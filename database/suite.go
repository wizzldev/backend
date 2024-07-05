package database

import (
	"fmt"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func MustConnectTestDB() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", configs.Env.Database.Username, configs.Env.Database.Password, configs.Env.Database.Host, configs.Env.Database.Port, configs.Env.Database.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err.Error())
	}

	log.Println("successfully connected to the database!")
	if configs.Env.Debug {
		db.Logger = logger.Default.LogMode(logger.Warn)
	} else {
		db.Logger = logger.Default.LogMode(logger.Error)
	}

	log.Println("Running migrations")
	err = db.AutoMigrate(getModels()...)

	if err != nil {
		log.Fatal("Failed to migrate: " + err.Error())
	}

	now := time.Now()
	pass, _ := utils.NewPassword("secret1234").Hash()
	// Creating main test user
	db.Create(&models.User{
		FirstName:       "Jane",
		LastName:        "Roe",
		Email:           "jane@example.com",
		Password:        pass, // password
		ImageURL:        "default.webp",
		EmailVerifiedAt: &now,
	})
	db.Create(&models.User{
		FirstName:       "Sam",
		LastName:        "Doe",
		Email:           "sam@example.com",
		Password:        pass, // password
		ImageURL:        "default.webp",
		EmailVerifiedAt: &now,
	})

	DB = db
}

func CleanUpTestDB() error {
	err := DB.Migrator().DropTable(getModels()...)
	if err != nil {
		return fmt.Errorf("failed to drop tables: " + err.Error())
	}

	DB.Rollback()

	db, err := DB.DB()
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}
