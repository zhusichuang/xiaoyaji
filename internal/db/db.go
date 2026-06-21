package db

import (
	"fmt"
	"time"

	"wxcloudrun-golang/internal/config"
	"wxcloudrun-golang/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var instance *gorm.DB

func Init() error {
	cfg := config.LoadMySQLConfig()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&readTimeout=1500ms&writeTimeout=1500ms",
		cfg.Username,
		cfg.Password,
		cfg.Address,
		cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(
		&model.User{},
		&model.Family{},
		&model.FamilyMember{},
		&model.Baby{},
		&model.BabyAction{},
	); err != nil {
		return err
	}

	instance = db
	return nil
}

func Get() *gorm.DB {
	return instance
}
