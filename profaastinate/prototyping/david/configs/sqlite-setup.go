package configs

import (
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupSqlite sets up the sqlite client based on the environment variables
func SetupSqlite(ctx context.Context) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../../../../../../docker/nuclio/sqlite/benchmark.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
