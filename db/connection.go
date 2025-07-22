package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "host=aws-0-sa-east-1.pooler.supabase.com user=postgres.ocucdyqtwdrbjkujmqbd password= dbname=postgres port=5432 sslmode=require"
var DB *gorm.DB

func DBConnection() {
	var error error
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil{
		log.Fatal(error)
	} else {
		log.Println("DB connected")
	}
}