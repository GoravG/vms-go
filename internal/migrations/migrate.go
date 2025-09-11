package migrations

import (
	"database/sql"
	"vms_go/internal/models"
	"vms_go/internal/utils"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) RunMigrtations() error {
	utils.LogInfof("Running Database Migrations")
	migrations := []struct {
		tablename string
		sql       string
	}{
		{"users", models.User{}.CreateTableSQL()},
	}

	for _, migration := range migrations {
		utils.LogInfof("Creating table: %s", migration.tablename)
		_, err := m.db.Exec(migration.sql)
		if err != nil {
			return err
		}
		utils.LogInfof("Table %s created successfully", migration.tablename)
	}
	utils.LogInfof("All migrations completed successfully")
	return nil
}
