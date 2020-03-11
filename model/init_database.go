package model

import (
	"fmt"
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	config "github.com/KevinAntonius/RegistrationCodingTest/config"
)

func Initialize() {
    string_param := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_SSLMODE)
	db, err := gorm.Open(config.DB_ADAPTER, string_param)
	if err != nil{
		fmt.Println("Connection Not Established")
		panic(err)
	}
	defer db.Close()
    
    db.Debug().AutoMigrate(&User{})
}