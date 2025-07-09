package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"webenable-cms-backend/utils"

	"github.com/go-kivik/kivik/v4"
	"github.com/sirupsen/logrus"
)

// BackupManager handles database backup and restore operations
type BackupManager struct {
	db        *OptimizedDB
	backupDir string
}

// BackupConfig contains backup configuration
type BackupConfig struct {
	BackupDir       string        `json:"backup_dir"`
	RetentionDays   int           `json:"retention_days"`
	BackupInterval  time.Duration `json:"backup_interval"`
	CompressBackups bool          `json:"compress_backups"`
}

// BackupMetadata contains information about a backup
type BackupMetadata struct {
	Timestamp  time.Time      `json:"timestamp"`
	Databases  []string       `json:"databases"`
	Size       int64          `json:"size"`
	Compressed bool           `json:"compressed"`
	Checksum   string         `json:"checksum"`
	Version    string         `json:"version"`
	DocCounts  map[string]int `json:"doc_counts"`
	Duration   time.Duration  `json:"duration"`
	Status     string         `json:"status"`
	Error      string         `json:"error,omitempty"`
}

// NewBackupManager creates a new backup manager
func NewBackupManager(db *OptimizedDB, backupDir string) *BackupManager {
	return &BackupManager{
		db:        db,
		backupDir: backupDir,
	}
}

// CreateBackup creates a full backup of all databases
func (bm *BackupManager) CreateBackup(ctx context.Context) (*BackupMetadata, error) {
	start := time.Now()
	timestamp := start.Format("2006-01-02_15-04-05")

	utils.LogInfo("Starting database backup", logrus.Fields{
		"timestamp": timestamp,
	})

	// Create backup directory
	backupPath := filepath.Join(bm.backupDir, fmt.Sprintf("backup_%s", timestamp))
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	metadata := &BackupMetadata{
		Timestamp: start,
		Databases: []string{"posts", "users", "contacts", "migrations"},
		DocCounts: make(map[string]int),
		Status:    "running",
	}

	var totalSize int64

	// Backup each database
	for _, dbName := range metadata.Databases {
		utils.LogInfo("Backing up database", logrus.Fields{
			"database": dbName,
		})

		size, count, err := bm.backupDatabase(ctx, dbName, backupPath)
		if err != nil {
			metadata.Status = "failed"
			metadata.Error = err.Error()
			return metadata, fmt.Errorf("failed to backup database %s: %w", dbName, err)
		}

		totalSize += size
		metadata.DocCounts[dbName] = count

		utils.LogInfo("Database backup completed", logrus.Fields{
			"database":  dbName,
			"doc_count": count,
			"size":      size,
		})
	}

	// Create metadata file
	metadata.Size = totalSize
	metadata.Duration = time.Since(start)
	metadata.Status = "completed"
	metadata.Version = "1.0"

	metadataFile := filepath.Join(backupPath, "metadata.json")
	if err := bm.saveMetadata(metadata, metadataFile); err != nil {
		return metadata, fmt.Errorf("failed to save metadata: %w", err)
	}

	utils.LogInfo("Database backup completed", logrus.Fields{
		"timestamp":   timestamp,
		"total_size":  totalSize,
		"duration":    metadata.Duration.String(),
		"total_docs":  bm.getTotalDocs(metadata.DocCounts),
		"backup_path": backupPath,
	})

	return metadata, nil
}

// backupDatabase backs up a single database
func (bm *BackupManager) backupDatabase(ctx context.Context, dbName, backupPath string) (int64, int, error) {
	db := bm.db.Client.DB(dbName)

	// Check if database exists
	if exists, _ := bm.db.Client.DBExists(ctx, dbName); !exists {
		utils.LogInfo("Database does not exist, skipping", logrus.Fields{
			"database": dbName,
		})
		return 0, 0, nil
	}

	// Create backup file
	backupFile := filepath.Join(backupPath, fmt.Sprintf("%s.json", dbName))
	file, err := os.Create(backupFile)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Get all documents
	rows := db.AllDocs(ctx, kivik.Param("include_docs", true))
	defer rows.Close()

	var documents []map[string]interface{}
	docCount := 0

	for rows.Next() {
		var doc map[string]interface{}
		if err := rows.ScanDoc(&doc); err != nil {
			utils.LogError(err, "Failed to scan document during backup", logrus.Fields{
				"database": dbName,
			})
			continue
		}

		// Add document ID and revision
		if id, err := rows.ID(); err == nil && id != "" {
			doc["_id"] = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			doc["_rev"] = rev
		}

		documents = append(documents, doc)
		docCount++
	}

	// Write documents to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(map[string]interface{}{
		"database":  dbName,
		"timestamp": time.Now(),
		"documents": documents,
		"count":     docCount,
	}); err != nil {
		return 0, 0, fmt.Errorf("failed to write backup data: %w", err)
	}

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return fileInfo.Size(), docCount, nil
}

// RestoreBackup restores a backup from the specified path
func (bm *BackupManager) RestoreBackup(ctx context.Context, backupPath string) error {
	utils.LogInfo("Starting database restore", logrus.Fields{
		"backup_path": backupPath,
	})

	// Load metadata
	metadataFile := filepath.Join(backupPath, "metadata.json")
	metadata, err := bm.loadMetadata(metadataFile)
	if err != nil {
		return fmt.Errorf("failed to load backup metadata: %w", err)
	}

	// Restore each database
	for _, dbName := range metadata.Databases {
		utils.LogInfo("Restoring database", logrus.Fields{
			"database": dbName,
		})

		if err := bm.restoreDatabase(ctx, dbName, backupPath); err != nil {
			return fmt.Errorf("failed to restore database %s: %w", dbName, err)
		}

		utils.LogInfo("Database restore completed", logrus.Fields{
			"database": dbName,
		})
	}

	utils.LogInfo("Database restore completed", logrus.Fields{
		"backup_path": backupPath,
		"databases":   len(metadata.Databases),
	})

	return nil
}

// restoreDatabase restores a single database
func (bm *BackupManager) restoreDatabase(ctx context.Context, dbName, backupPath string) error {
	backupFile := filepath.Join(backupPath, fmt.Sprintf("%s.json", dbName))

	// Check if backup file exists
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		utils.LogInfo("Backup file does not exist, skipping", logrus.Fields{
			"database": dbName,
			"file":     backupFile,
		})
		return nil
	}

	// Read backup file
	file, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	var backupData struct {
		Database  string                   `json:"database"`
		Timestamp time.Time                `json:"timestamp"`
		Documents []map[string]interface{} `json:"documents"`
		Count     int                      `json:"count"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&backupData); err != nil {
		return fmt.Errorf("failed to decode backup data: %w", err)
	}

	// Recreate database
	if exists, _ := bm.db.Client.DBExists(ctx, dbName); exists {
		if err := bm.db.Client.DestroyDB(ctx, dbName); err != nil {
			return fmt.Errorf("failed to destroy existing database: %w", err)
		}
	}

	if err := bm.db.Client.CreateDB(ctx, dbName); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Restore documents
	db := bm.db.Client.DB(dbName)

	// Prepare documents for bulk insert
	var docs []interface{}
	for _, doc := range backupData.Documents {
		// Remove revision for restore
		delete(doc, "_rev")
		docs = append(docs, doc)
	}

	// Bulk insert documents
	if len(docs) > 0 {
		results, err := db.BulkDocs(ctx, docs)
		if err != nil {
			return fmt.Errorf("failed to bulk insert documents: %w", err)
		}

		// Check for errors
		errorCount := 0
		for _, result := range results {
			if result.Error != nil {
				errorCount++
				utils.LogError(nil, "Document restore error", logrus.Fields{
					"database": dbName,
					"doc_id":   result.ID,
					"error":    result.Error.Error(),
				})
			}
		}

		if errorCount > 0 {
			utils.LogError(nil, "Some documents failed to restore", logrus.Fields{
				"database":    dbName,
				"error_count": errorCount,
				"total_docs":  len(docs),
			})
		}
	}

	return nil
}

// CleanupOldBackups removes backups older than the retention period
func (bm *BackupManager) CleanupOldBackups(retentionDays int) error {
	utils.LogInfo("Starting backup cleanup", logrus.Fields{
		"retention_days": retentionDays,
	})

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	removedCount := 0
	for _, entry := range entries {
		if !entry.IsDir() || !filepath.HasPrefix(entry.Name(), "backup_") {
			continue
		}

		backupPath := filepath.Join(bm.backupDir, entry.Name())
		metadataFile := filepath.Join(backupPath, "metadata.json")

		metadata, err := bm.loadMetadata(metadataFile)
		if err != nil {
			utils.LogError(err, "Failed to load backup metadata for cleanup", logrus.Fields{
				"backup_path": backupPath,
			})
			continue
		}

		if metadata.Timestamp.Before(cutoffTime) {
			if err := os.RemoveAll(backupPath); err != nil {
				utils.LogError(err, "Failed to remove old backup", logrus.Fields{
					"backup_path": backupPath,
				})
				continue
			}

			utils.LogInfo("Removed old backup", logrus.Fields{
				"backup_path": backupPath,
				"timestamp":   metadata.Timestamp,
			})
			removedCount++
		}
	}

	utils.LogInfo("Backup cleanup completed", logrus.Fields{
		"removed_count": removedCount,
	})

	return nil
}

// ListBackups returns a list of available backups
func (bm *BackupManager) ListBackups() ([]*BackupMetadata, error) {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []*BackupMetadata
	for _, entry := range entries {
		if !entry.IsDir() || !filepath.HasPrefix(entry.Name(), "backup_") {
			continue
		}

		backupPath := filepath.Join(bm.backupDir, entry.Name())
		metadataFile := filepath.Join(backupPath, "metadata.json")

		metadata, err := bm.loadMetadata(metadataFile)
		if err != nil {
			utils.LogError(err, "Failed to load backup metadata", logrus.Fields{
				"backup_path": backupPath,
			})
			continue
		}

		backups = append(backups, metadata)
	}

	return backups, nil
}

// saveMetadata saves backup metadata to a file
func (bm *BackupManager) saveMetadata(metadata *BackupMetadata, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metadata)
}

// loadMetadata loads backup metadata from a file
func (bm *BackupManager) loadMetadata(filename string) (*BackupMetadata, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata BackupMetadata
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// getTotalDocs calculates total document count across all databases
func (bm *BackupManager) getTotalDocs(docCounts map[string]int) int {
	total := 0
	for _, count := range docCounts {
		total += count
	}
	return total
}

// ScheduledBackupService provides automated backup scheduling
type ScheduledBackupService struct {
	backupManager *BackupManager
	config        BackupConfig
	stopChan      chan struct{}
}

// NewScheduledBackupService creates a new scheduled backup service
func NewScheduledBackupService(backupManager *BackupManager, config BackupConfig) *ScheduledBackupService {
	return &ScheduledBackupService{
		backupManager: backupManager,
		config:        config,
		stopChan:      make(chan struct{}),
	}
}

// Start starts the scheduled backup service
func (sbs *ScheduledBackupService) Start(ctx context.Context) {
	ticker := time.NewTicker(sbs.config.BackupInterval)
	defer ticker.Stop()

	utils.LogInfo("Scheduled backup service started", logrus.Fields{
		"interval":       sbs.config.BackupInterval.String(),
		"retention_days": sbs.config.RetentionDays,
	})

	for {
		select {
		case <-ticker.C:
			// Create backup
			if _, err := sbs.backupManager.CreateBackup(ctx); err != nil {
				utils.LogError(err, "Scheduled backup failed", logrus.Fields{})
			}

			// Cleanup old backups
			if err := sbs.backupManager.CleanupOldBackups(sbs.config.RetentionDays); err != nil {
				utils.LogError(err, "Backup cleanup failed", logrus.Fields{})
			}

		case <-sbs.stopChan:
			utils.LogInfo("Scheduled backup service stopped", logrus.Fields{})
			return

		case <-ctx.Done():
			utils.LogInfo("Scheduled backup service stopped due to context cancellation", logrus.Fields{})
			return
		}
	}
}

// Stop stops the scheduled backup service
func (sbs *ScheduledBackupService) Stop() {
	close(sbs.stopChan)
}
