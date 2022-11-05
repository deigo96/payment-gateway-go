package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseDriver string

const (
	PostgreSQL DatabaseDriver = "postgres"
)

type DatabaseConnection struct {
	Driver     DatabaseDriver
	PostgreSQL *gorm.DB
}

func NewDatabaseConnection(config *AppConfig) *DatabaseConnection {
	var db DatabaseConnection

	if config.Driver != "PostgreSQL" {
		panic("Database driver not supported")
	}

	db.Driver = PostgreSQL
	db.PostgreSQL = NewPostgreSQL(config)

	return &db
}

func NewPostgreSQL(config *AppConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DB_Host,
		config.DB_User,
		config.DB_Pass,
		config.DB_Name,
		config.DB_Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	return db
}

func (db *DatabaseConnection) CloseConnection() {
	if db.PostgreSQL != nil {
		db, _ := db.PostgreSQL.DB()
		db.Close()
	}
}
