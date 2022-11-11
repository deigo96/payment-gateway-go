package config

import (
	_ "github.com/denisenkom/go-mssqldb"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseDriver string

const (
	SqlServer DatabaseDriver = "sqlserver"
)

type DatabaseConnection struct {
	Driver    DatabaseDriver
	SqlServer *gorm.DB
}

func NewDatabaseConnection(config *AppConfig) *DatabaseConnection {
	var db DatabaseConnection

	if config.DRIVER != "SqlServer" {
		panic("Database driver not supported")
	}

	db.Driver = SqlServer
	db.SqlServer = NewSqlServer(config)

	return &db
}

func NewSqlServer(config *AppConfig) *gorm.DB {
	dsn := "server=" + config.DB_HOST + ";user id=" + config.DB_USER + ";password=" + config.DB_PASS + ";database=" + config.DB_NAME + ";encrypt=disable;"

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
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
	if db.SqlServer != nil {
		db, _ := db.SqlServer.DB()
		db.Close()
	}
}
