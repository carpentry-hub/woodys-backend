package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "postgresql://postgres.ocucdyqtwdrbjkujmqbd:[YOUR-PASSWORD]@aws-0-sa-east-1.pooler.supabase.com:5432/postgres"
var DB *gorm.DB

func DBConnection() {
	var error error
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil{
		log.Fatal(error)
	} else{
		log.Println("DB connected")
	}
}