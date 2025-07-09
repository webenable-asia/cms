package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"webenable-cms-backend/utils"

	"github.com/sirupsen/logrus"
)

// DatabaseMonitor provides comprehensive database monitoring
type DatabaseMonitor struct {
	db              *OptimizedDB
	metrics         *DatabaseMetrics
	healthChecks    []HealthCheck
	alertThresholds AlertThresholds
	mutex           sync.RWMutex
}

// DatabaseMetrics contains database performance metrics
type DatabaseMetrics struct {
	ConnectionPool ConnectionPoolMetrics    `json:"connection_pool"`
	QueryMetrics   map[string]*QueryMetrics `json:"query_metrics"`
	DatabaseSizes  map[string]int64         `json:"database_sizes"`
	DocumentCounts map[string]int           `json:"document_counts"`
	IndexMetrics   map[string]IndexMetrics  `json:"index_metrics"`
	HealthStatus   HealthStatus             `json:"health_status"`
	LastUpdated    time.Time                `json:"last_updated"`
}

// ConnectionPoolMetrics tracks connection pool performance
type ConnectionPoolMetrics struct {
	MaxConnections    int           `json:"max_connections"`
	ActiveConnections int           `json:"active_connections"`
	IdleConnections   int           `json:"idle_connections"`
	WaitingRequests   int           `json:"waiting_requests"`
	AverageWaitTime   time.Duration `json:"average_wait_time"`
	TotalRequests     int64         `json:"total_requests"`
}

// IndexMetrics tracks index usage and performance
type IndexMetrics struct {
	Name        string        `json:"name"`
	Database    string        `json:"database"`
	Fields      []string      `json:"fields"`
	UsageCount  int64         `json:"usage_count"`
	LastUsed    time.Time     `json:"last_used"`
	AverageTime time.Duration `json:"average_time"`
	Size        int64         `json:"size"`
}

// HealthStatus represents overall database health
type HealthStatus struct {
	Status       string            `json:"status"` // healthy, warning, critical
	LastCheck    time.Time         `json:"last_check"`
	Issues       []HealthIssue     `json:"issues"`
	Uptime       time.Duration     `json:"uptime"`
	ResponseTime time.Duration     `json:"response_time"`
	Databases    map[string]string `json:"databases"` // database -> status
}

// HealthIssue represents a health check issue
type HealthIssue struct {
	Type      string    `json:"type"`
	Severity  string    `json:"severity"` // info, warning, critical
	Message   string    `json:"message"`
	Database  string    `json:"database,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Resolved  bool      `json:"resolved"`
}

// AlertThresholds defines thresholds for alerts
type AlertThresholds struct {
	MaxResponseTime    time.Duration `json:"max_response_time"`
	MaxConnectionUsage float64       `json:"max_connection_usage"`
	MaxQueryTime       time.Duration `json:"max_query_time"`
	MinDiskSpace       int64         `json:"min_disk_space"`
	MaxErrorRate       float64       `json:"max_error_rate"`
}

// HealthCheck represents a health check function
type HealthCheck struct {
	Name        string
	Description string
	Check       func(ctx context.Context, db *OptimizedDB) *HealthIssue
	Interval    time.Duration
	LastRun     time.Time
}

// NewDatabaseMonitor creates a new database monitor
func NewDatabaseMonitor(db *OptimizedDB) *DatabaseMonitor {
	return &DatabaseMonitor{
		db: db,
		metrics: &DatabaseMetrics{
			QueryMetrics:   make(map[string]*QueryMetrics),
			DatabaseSizes:  make(map[string]int64),
			DocumentCounts: make(map[string]int),
			IndexMetrics:   make(map[string]IndexMetrics),
			HealthStatus: HealthStatus{
				Status:    "unknown",
				Databases: make(map[string]string),
			},
		},
		healthChecks: getDefaultHealthChecks(),
		alertThresholds: AlertThresholds{
			MaxResponseTime:    5 * time.Second,
			MaxConnectionUsage: 0.8, // 80%
			MaxQueryTime:       10 * time.Second,
			MinDiskSpace:       1024 * 1024 * 1024, // 1GB
			MaxErrorRate:       0.05,               // 5%
		},
	}
}

// getDefaultHealthChecks returns the default set of health checks
func getDefaultHealthChecks() []HealthCheck {
	return []HealthCheck{
		{
			Name:        "database_connectivity",
			Description: "Check database connectivity and response time",
			Interval:    30 * time.Second,
			Check: func(ctx context.Context, db *OptimizedDB) *HealthIssue {
				start := time.Now()
				if err := db.Health(); err != nil {
					return &HealthIssue{
						Type:      "connectivity",
						Severity:  "critical",
						Message:   fmt.Sprintf("Database connectivity failed: %v", err),
						Timestamp: time.Now(),
					}
				}

				responseTime := time.Since(start)
				if responseTime > 5*time.Second {
					return &HealthIssue{
						Type:      "performance",
						Severity:  "warning",
						Message:   fmt.Sprintf("Database response time is slow: %v", responseTime),
						Timestamp: time.Now(),
					}
				}

				return nil
			},
		},
		{
			Name:        "connection_pool_usage",
			Description: "Monitor connection pool usage",
			Interval:    1 * time.Minute,
			Check: func(ctx context.Context, db *OptimizedDB) *HealthIssue {
				activeConnections := len(db.connPool)
				usage := float64(activeConnections) / float64(db.maxConnections)

				if usage > 0.9 {
					return &HealthIssue{
						Type:      "resource",
						Severity:  "warning",
						Message:   fmt.Sprintf("Connection pool usage is high: %.1f%%", usage*100),
						Timestamp: time.Now(),
					}
				}

				return nil
			},
		},
		{
			Name:        "database_sizes",
			Description: "Monitor database sizes and growth",
			Interval:    5 * time.Minute,
			Check: func(ctx context.Context, db *OptimizedDB) *HealthIssue {
				databases := []string{"posts", "users", "contacts", "migrations"}

				for _, dbName := range databases {
					// This is a simplified check - in production you'd want to
					// implement actual size monitoring using CouchDB's _db_info endpoint
					if dbExists, _ := db.Client.DBExists(ctx, dbName); !dbExists {
						return &HealthIssue{
							Type:      "configuration",
							Severity:  "critical",
							Message:   fmt.Sprintf("Database %s does not exist", dbName),
							Database:  dbName,
							Timestamp: time.Now(),
						}
					}
				}

				return nil
			},
		},
		{
			Name:        "slow_queries",
			Description: "Detect slow-running queries",
			Interval:    2 * time.Minute,
			Check: func(ctx context.Context, db *OptimizedDB) *HealthIssue {
				metrics := db.GetQueryMetrics()

				for queryName, metric := range metrics {
					if metric.AverageTime > 10*time.Second {
						return &HealthIssue{
							Type:      "performance",
							Severity:  "warning",
							Message:   fmt.Sprintf("Query %s has slow average time: %v", queryName, metric.AverageTime),
							Timestamp: time.Now(),
						}
					}
				}

				return nil
			},
		},
	}
}

// StartMonitoring starts the database monitoring service
func (dm *DatabaseMonitor) StartMonitoring(ctx context.Context) {
	utils.LogInfo("Starting database monitoring", logrus.Fields{
		"health_checks": len(dm.healthChecks),
	})

	// Start health check routines
	for _, healthCheck := range dm.healthChecks {
		go dm.runHealthCheck(ctx, healthCheck)
	}

	// Start metrics collection
	go dm.collectMetrics(ctx)

	// Start periodic reporting
	go dm.periodicReport(ctx)
}

// runHealthCheck runs a single health check periodically
func (dm *DatabaseMonitor) runHealthCheck(ctx context.Context, healthCheck HealthCheck) {
	ticker := time.NewTicker(healthCheck.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			issue := healthCheck.Check(ctx, dm.db)
			if issue != nil {
				dm.recordHealthIssue(*issue)
				utils.LogError(nil, "Health check failed", logrus.Fields{
					"check":    healthCheck.Name,
					"severity": issue.Severity,
					"message":  issue.Message,
				})
			}

			dm.mutex.Lock()
			healthCheck.LastRun = time.Now()
			dm.mutex.Unlock()

		case <-ctx.Done():
			return
		}
	}
}

// collectMetrics collects database metrics periodically
func (dm *DatabaseMonitor) collectMetrics(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dm.updateMetrics(ctx)

		case <-ctx.Done():
			return
		}
	}
}

// updateMetrics updates all database metrics
func (dm *DatabaseMonitor) updateMetrics(ctx context.Context) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Update connection pool metrics
	dm.metrics.ConnectionPool = ConnectionPoolMetrics{
		MaxConnections:    dm.db.maxConnections,
		ActiveConnections: dm.db.maxConnections - len(dm.db.connPool),
		IdleConnections:   len(dm.db.connPool),
	}

	// Update query metrics
	dm.metrics.QueryMetrics = dm.db.GetQueryMetrics()

	// Update document counts
	databases := []string{"posts", "users", "contacts", "migrations"}
	for _, dbName := range databases {
		if count, err := dm.getDocumentCount(ctx, dbName); err == nil {
			dm.metrics.DocumentCounts[dbName] = count
		}
	}

	// Update health status
	dm.updateHealthStatus()

	dm.metrics.LastUpdated = time.Now()
}

// getDocumentCount gets the document count for a database
func (dm *DatabaseMonitor) getDocumentCount(ctx context.Context, dbName string) (int, error) {
	if exists, _ := dm.db.Client.DBExists(ctx, dbName); !exists {
		return 0, nil
	}

	db := dm.db.Client.DB(dbName)
	rows := db.AllDocs(ctx)
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

// updateHealthStatus updates the overall health status
func (dm *DatabaseMonitor) updateHealthStatus() {
	status := "healthy"
	criticalIssues := 0
	warningIssues := 0

	for _, issue := range dm.metrics.HealthStatus.Issues {
		if !issue.Resolved {
			switch issue.Severity {
			case "critical":
				criticalIssues++
			case "warning":
				warningIssues++
			}
		}
	}

	if criticalIssues > 0 {
		status = "critical"
	} else if warningIssues > 0 {
		status = "warning"
	}

	dm.metrics.HealthStatus.Status = status
	dm.metrics.HealthStatus.LastCheck = time.Now()
}

// recordHealthIssue records a health issue
func (dm *DatabaseMonitor) recordHealthIssue(issue HealthIssue) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Check if this issue already exists and is unresolved
	for i, existingIssue := range dm.metrics.HealthStatus.Issues {
		if existingIssue.Type == issue.Type &&
			existingIssue.Database == issue.Database &&
			!existingIssue.Resolved {
			// Update existing issue
			dm.metrics.HealthStatus.Issues[i] = issue
			return
		}
	}

	// Add new issue
	dm.metrics.HealthStatus.Issues = append(dm.metrics.HealthStatus.Issues, issue)

	// Keep only the last 100 issues
	if len(dm.metrics.HealthStatus.Issues) > 100 {
		dm.metrics.HealthStatus.Issues = dm.metrics.HealthStatus.Issues[len(dm.metrics.HealthStatus.Issues)-100:]
	}
}

// periodicReport generates periodic monitoring reports
func (dm *DatabaseMonitor) periodicReport(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dm.generateReport()

		case <-ctx.Done():
			return
		}
	}
}

// generateReport generates a monitoring report
func (dm *DatabaseMonitor) generateReport() {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	utils.LogInfo("Database monitoring report", logrus.Fields{
		"status":             dm.metrics.HealthStatus.Status,
		"active_connections": dm.metrics.ConnectionPool.ActiveConnections,
		"total_documents":    dm.getTotalDocuments(),
		"query_count":        len(dm.metrics.QueryMetrics),
		"unresolved_issues":  dm.getUnresolvedIssueCount(),
	})

	// Log slow queries
	for queryName, metric := range dm.metrics.QueryMetrics {
		if metric.AverageTime > 5*time.Second {
			utils.LogError(nil, "Slow query detected", logrus.Fields{
				"query":        queryName,
				"average_time": metric.AverageTime.String(),
				"count":        metric.Count,
			})
		}
	}
}

// getTotalDocuments returns the total number of documents across all databases
func (dm *DatabaseMonitor) getTotalDocuments() int {
	total := 0
	for _, count := range dm.metrics.DocumentCounts {
		total += count
	}
	return total
}

// getUnresolvedIssueCount returns the number of unresolved issues
func (dm *DatabaseMonitor) getUnresolvedIssueCount() int {
	count := 0
	for _, issue := range dm.metrics.HealthStatus.Issues {
		if !issue.Resolved {
			count++
		}
	}
	return count
}

// GetMetrics returns the current database metrics
func (dm *DatabaseMonitor) GetMetrics() *DatabaseMetrics {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	// Return a copy to avoid race conditions
	metricsCopy := *dm.metrics
	return &metricsCopy
}

// GetHealthStatus returns the current health status
func (dm *DatabaseMonitor) GetHealthStatus() HealthStatus {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	return dm.metrics.HealthStatus
}

// ResolveIssue marks a health issue as resolved
func (dm *DatabaseMonitor) ResolveIssue(issueType, database string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for i, issue := range dm.metrics.HealthStatus.Issues {
		if issue.Type == issueType && issue.Database == database && !issue.Resolved {
			dm.metrics.HealthStatus.Issues[i].Resolved = true
			utils.LogInfo("Health issue resolved", logrus.Fields{
				"type":     issueType,
				"database": database,
			})
			break
		}
	}
}
