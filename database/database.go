package database

import (
	"fmt"
	"log"

	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func MustConnect() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", configs.Env.Database.Username, configs.Env.Database.Password, configs.Env.Database.Host, configs.Env.Database.Port, configs.Env.Database.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to the database: %v\n", err.Error()))
	}

	log.Println("successfully connected to the database!")
	if configs.Env.Debug {
		db.Logger = logger.Default.LogMode(logger.Info)
	} else {
		db.Logger = logger.Default.LogMode(logger.Error)
	}

	log.Println("Running migrations")
	err = db.AutoMigrate(
		&models.Message{},
		&models.MessageLike{},
		&models.Group{},
		&models.MemberRole{},
		&models.Block{},
		&models.EmailVerification{},
		&models.ResetPassword{},
		&models.Bot{},
		&models.GroupUser{},
		&models.User{},
		&models.File{},
	)

	if err != nil {
		log.Fatal("Failed to migrate: " + err.Error())
	}

	DB = db
}
