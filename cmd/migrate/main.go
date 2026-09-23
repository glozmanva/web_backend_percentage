package main

import (
	"log"

	"deposit_month/internal/app/ds"
	"deposit_month/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	db, err := gorm.Open(
		postgres.Open(dsn.FromEnv()),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.DepositMonth{},
		&ds.DepositMonthLike{},
	)
	if err != nil {
		log.Fatal("cant migrate db: ", err)
	}

	err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS one_draft_per_user
		ON deposit_months (creator_id)
		WHERE status = 'draft'
	`).Error
	if err != nil {
		log.Fatal("cant create draft index: ", err)
	}

	log.Println("migration completed")
}
