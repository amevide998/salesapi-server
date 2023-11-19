package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"sales-api/Model"
)

var DB *gorm.DB

func Connect() {
	// db
	err := godotenv.Load()
	if err != nil {
		return
	}
	dsn := os.Getenv("DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("db connection failed")
	}

	DB = db
	fmt.Println("db connection successfull")
	//AutoMigrate(db)
}

func AutoMigrate(c *gorm.DB) {
	err := c.Debug().AutoMigrate(
		&Model.Cashier{},
		&Model.Category{},
		&Model.Payment{},
		&Model.PaymentType{},
		&Model.Product{},
		&Model.Discount{},
		&Model.Order{},
	)
	if err != nil {
		return
	}
}
