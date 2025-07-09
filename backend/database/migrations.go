package database

import (
	"context"
	"fmt"
	"sort"
	"time"

	"webenable-cms-backend/utils"

	"github.com/sirupsen/logrus"
)

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	Up          func(ctx context.Context, db *OptimizedDB) error
	Down        func(ctx context.Context, db *OptimizedDB) error
}

// MigrationRecord tracks applied migrations
type MigrationRecord struct {
	ID          string    `json:"_id,omitempty"`
	Rev         string    `json:"_rev,omitempty"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	AppliedAt   time.Time `json:"applied_at"`
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *OptimizedDB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *OptimizedDB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: getMigrations(),
	}
}

// getMigrations returns all available migrations in order
func getMigrations() []Migration {
	return []Migration{
		{
			Version:     "001_initial_indexes",
			Description: "Create initial database indexes for performance",
			Up: func(ctx context.Context, db *OptimizedDB) error {
				return db.createIndexes()
			},
			Down: func(ctx context.Context, db *OptimizedDB) error {
				// Drop indexes (implementation depends on requirements)
				return nil
			},
		},
		{
			Version:     "002_add_view_count_index",
			Description: "Add index for post view count sorting",
			Up: func(ctx context.Context, db *OptimizedDB) error {
				index := map[string]interface{}{
					"fields": []string{"view_count", "status"},
				}
				return db.PostsDB.CreateIndex(ctx, "", "view-count-status-index", index)
			},
			Down: func(ctx context.Context, db *OptimizedDB) error {
				// Drop the index
				return nil
			},
		},
		{
			Version:     "003_add_reading_time_index",
			Description: "Add index for reading time filtering",
			Up: func(ctx context.Context, db *OptimizedDB) error {
				index := map[string]interface{}{
					"fields": []string{"reading_time", "status"},
				}
				return db.PostsDB.CreateIndex(ctx, "", "reading-time-status-index", index)
			},
			Down: func(ctx context.Context, db *OptimizedDB) error {
				// Drop the index
				return nil
			},
		},
		{
			Version:     "004_optimize_contact_queries",
			Description: "Add compound indexes for contact management",
			Up: func(ctx context.Context, db *OptimizedDB) error {
				// Add index for email + status queries
				emailStatusIndex := map[string]interface{}{
					"fields": []string{"email", "status"},
				}
				if err := db.ContactsDB.CreateIndex(ctx, "", "email-status-index", emailStatusIndex); err != nil {
					return err
				}

				// Add index for company queries
				companyIndex := map[string]interface{}{
					"fields": []string{"company", "created_at"},
				}
				return db.ContactsDB.CreateIndex(ctx, "", "company-created-index", companyIndex)
			},
			Down: func(ctx context.Context, db *OptimizedDB) error {
				// Drop the indexes
				return nil
			},
		},
	}
}

// setupMigrationsDB creates the migrations tracking database
func (m *MigrationManager) setupMigrationsDB(ctx context.Context) error {
	if exists, _ := m.db.Client.DBExists(ctx, "migrations"); !exists {
		if err := m.db.Client.CreateDB(ctx, "migrations"); err != nil {
			return fmt.Errorf("failed to create migrations database: %w", err)
		}
		utils.LogInfo("Created migrations database", logrus.Fields{})
	}
	return nil
}

// getAppliedMigrations returns a list of applied migration versions
func (m *MigrationManager) getAppliedMigrations(ctx context.Context) (map[string]MigrationRecord, error) {
	migrationsDB := m.db.Client.DB("migrations")

	query := map[string]interface{}{
		"selector": map[string]interface{}{},
	}

	rows := migrationsDB.Find(ctx, query)
	defer rows.Close()

	applied := make(map[string]MigrationRecord)
	for rows.Next() {
		var record MigrationRecord
		if err := rows.ScanDoc(&record); err != nil {
			continue
		}

		if id, err := rows.ID(); err == nil && id != "" {
			record.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			record.Rev = rev
		}

		applied[record.Version] = record
	}

	return applied, nil
}

// recordMigration records a successful migration
func (m *MigrationManager) recordMigration(ctx context.Context, migration Migration) error {
	migrationsDB := m.db.Client.DB("migrations")

	record := MigrationRecord{
		ID:          migration.Version,
		Version:     migration.Version,
		Description: migration.Description,
		AppliedAt:   time.Now(),
	}

	_, err := migrationsDB.Put(ctx, migration.Version, record)
	if err != nil {
		return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
	}

	return nil
}

// removeMigrationRecord removes a migration record (for rollback)
func (m *MigrationManager) removeMigrationRecord(ctx context.Context, version string) error {
	migrationsDB := m.db.Client.DB("migrations")

	// Get the record to get its revision
	var record MigrationRecord
	row := migrationsDB.Get(ctx, version)
	if err := row.ScanDoc(&record); err != nil {
		return fmt.Errorf("failed to get migration record %s: %w", version, err)
	}

	if rev, err := row.Rev(); err == nil && rev != "" {
		record.Rev = rev
	}

	_, err := migrationsDB.Delete(ctx, version, record.Rev)
	if err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", version, err)
	}

	return nil
}

// Migrate runs all pending migrations
func (m *MigrationManager) Migrate(ctx context.Context) error {
	// Setup migrations database
	if err := m.setupMigrationsDB(ctx); err != nil {
		return err
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	// Run pending migrations
	for _, migration := range m.migrations {
		if _, exists := applied[migration.Version]; exists {
			utils.LogInfo("Migration already applied", logrus.Fields{
				"version": migration.Version,
			})
			continue
		}

		utils.LogInfo("Applying migration", logrus.Fields{
			"version":     migration.Version,
			"description": migration.Description,
		})

		start := time.Now()
		if err := migration.Up(ctx, m.db); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.Version, err)
		}

		// Record successful migration
		if err := m.recordMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		utils.LogInfo("Migration completed", logrus.Fields{
			"version":  migration.Version,
			"duration": time.Since(start).String(),
		})
	}

	utils.LogInfo("All migrations completed", logrus.Fields{
		"total_migrations": len(m.migrations),
	})

	return nil
}

// Rollback rolls back the last migration
func (m *MigrationManager) Rollback(ctx context.Context) error {
	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		utils.LogInfo("No migrations to rollback", logrus.Fields{})
		return nil
	}

	// Find the latest migration
	var latestMigration *Migration
	var latestVersion string
	var latestTime time.Time

	for version, record := range applied {
		if record.AppliedAt.After(latestTime) {
			latestTime = record.AppliedAt
			latestVersion = version

			// Find the migration definition
			for _, migration := range m.migrations {
				if migration.Version == version {
					latestMigration = &migration
					break
				}
			}
		}
	}

	if latestMigration == nil {
		return fmt.Errorf("migration definition not found for version %s", latestVersion)
	}

	utils.LogInfo("Rolling back migration", logrus.Fields{
		"version":     latestMigration.Version,
		"description": latestMigration.Description,
	})

	start := time.Now()
	if err := latestMigration.Down(ctx, m.db); err != nil {
		return fmt.Errorf("rollback of migration %s failed: %w", latestMigration.Version, err)
	}

	// Remove migration record
	if err := m.removeMigrationRecord(ctx, latestMigration.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", latestMigration.Version, err)
	}

	utils.LogInfo("Migration rolled back", logrus.Fields{
		"version":  latestMigration.Version,
		"duration": time.Since(start).String(),
	})

	return nil
}

// Status returns the current migration status
func (m *MigrationManager) Status(ctx context.Context) ([]MigrationStatus, error) {
	// Setup migrations database
	if err := m.setupMigrationsDB(ctx); err != nil {
		return nil, err
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	var status []MigrationStatus
	for _, migration := range m.migrations {
		s := MigrationStatus{
			Version:     migration.Version,
			Description: migration.Description,
			Applied:     false,
		}

		if record, exists := applied[migration.Version]; exists {
			s.Applied = true
			s.AppliedAt = &record.AppliedAt
		}

		status = append(status, s)
	}

	return status, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Applied     bool       `json:"applied"`
	AppliedAt   *time.Time `json:"applied_at,omitempty"`
}
