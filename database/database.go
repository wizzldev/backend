package database

import (
	"fmt"
	"github.com/wizzldev/chat/database/models"
	"log"

	"github.com/wizzldev/chat/pkg/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func getModels() []interface{} {
	return []interface{}{
		&models.Message{},
		&models.MessageLike{},
		&models.Group{},
		&models.Block{},
		&models.EmailVerification{},
		&models.ResetPassword{},
		&models.Theme{},
		&models.GroupUser{},
		&models.AllowedIP{},
		&models.Session{},
		&models.User{},
		&models.File{},
	}
}

func MustConnect() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", configs.Env.Database.Username, configs.Env.Database.Password, configs.Env.Database.Host, configs.Env.Database.Port, configs.Env.Database.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true,
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

	DB = db
}
