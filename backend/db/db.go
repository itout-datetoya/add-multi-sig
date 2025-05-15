package db

import (
	"os"
	"fmt"
	"log"

	"multisigservice/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB はPostgreSQLへの接続を確立し、必要なテーブルを自動マイグレーションします。
func InitDB() {
	dsn := os.Getenv("DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// モデルのスキーマを自動作成／更新
	err = DB.AutoMigrate(&models.User{}, &models.MultiSig{})
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	fmt.Println("Database connection established and schema migrated.")
}
