package infrastructure

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/philaden/mds-stock-keeping/application/domains"
)

var Connection *gorm.DB

func SetUpDatabaseServices(dbConfig AppConfiguration) {

	connectionString := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbUser, dbConfig.DbName, dbConfig.DbPassword)

	fmt.Printf("connection string %s", connectionString)

	db, err := gorm.Open(dbConfig.DbDialect, connectionString)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.LogMode(true)

	Connection = db

	dbase := db.DB()
	err = dbase.Ping()

	if err != nil {
		panic(err.Error())
	}

	if err := MigrateDatabase(); err != nil {
		fmt.Print("There was an error while trying to migrate tables..")
		os.Exit(1)
	}
}

func MigrateDatabase() error {

	if err := Connection.AutoMigrate(&domains.Product{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&domains.Order{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&domains.OrderItem{}).Error; err != nil {
		return err
	}

	return nil
}
