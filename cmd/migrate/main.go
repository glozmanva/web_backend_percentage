package main

import (
	"fmt"

	"deposit_month/internal/app/ds"
	"deposit_month/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func execOrPanic(db *gorm.DB, query string) {
	if err := db.Exec(query).Error; err != nil {
		panic(err)
	}
}

func main() {
	_ = godotenv.Load()

	db, err := gorm.Open(
		postgres.Open(dsn.FromEnv()),
		&gorm.Config{},
	)

	if err != nil {
		panic("failed to connect database")
	}

	if db.Migrator().HasTable(&ds.User{}) {
		execOrPanic(
			db,
			`
			ALTER TABLE users
			ADD COLUMN IF NOT EXISTS password VARCHAR(255)
			`,
		)

		execOrPanic(
			db,
			`
			UPDATE users
			SET password = 'not_set'
			WHERE password IS NULL
			   OR password = ''
			`,
		)

		execOrPanic(
			db,
			`
			ALTER TABLE users
			ALTER COLUMN password SET NOT NULL
			`,
		)
	}

	execOrPanic(
		db,
		`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name = 'deposit_months'
				  AND column_name = 'short_description'
			)
			AND EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name = 'deposit_months'
				  AND column_name = 'description'
			)
			THEN
				ALTER TABLE deposit_months
				ALTER COLUMN short_description TYPE VARCHAR(1000);

				UPDATE deposit_months
				SET short_description = description
				WHERE description IS NOT NULL;

				ALTER TABLE deposit_months
				DROP COLUMN description;

				ALTER TABLE deposit_months
				RENAME COLUMN short_description TO description;

			ELSIF EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name = 'deposit_months'
				  AND column_name = 'short_description'
			)
			AND NOT EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name = 'deposit_months'
				  AND column_name = 'description'
			)
			THEN
				ALTER TABLE deposit_months
				ALTER COLUMN short_description TYPE VARCHAR(1000);

				ALTER TABLE deposit_months
				RENAME COLUMN short_description TO description;
			END IF;
		END $$;
		`,
	)

	err = db.AutoMigrate(
		&ds.User{},
		&ds.DepositMonth{},
		&ds.DepositMonthLike{},
	)

	if err != nil {
		panic("cant migrate db")
	}

	execOrPanic(
		db,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS one_draft_per_user
		ON deposit_months (creator_id)
		WHERE status = 'draft'
		`,
	)

	fmt.Println("migration completed")
}
