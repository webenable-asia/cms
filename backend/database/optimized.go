package database

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"webenable-cms-backend/utils"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
	"github.com/sirupsen/logrus"
)

// OptimizedDB provides an optimized database layer with connection pooling,
// indexing, and query optimization
type OptimizedDB struct {
	Client     *kivik.Client
	PostsDB    *kivik.DB
	UsersDB    *kivik.DB
	ContactsDB *kivik.DB

	// Connection pool settings
	maxConnections int
	connPool       chan *kivik.Client
	connMutex      sync.RWMutex

	// Query cache
	queryCache map[string]interface{}
	cacheMutex sync.RWMutex

	// Performance metrics
	queryMetrics map[string]*QueryMetrics
	metricsMutex sync.RWMutex
}

type QueryMetrics struct {
	Count       int64
	TotalTime   time.Duration
	AverageTime time.Duration
	LastUsed    time.Time
}

var OptimizedInstance *OptimizedDB

// InitOptimized initializes the optimized database with connection pooling and indexing
func InitOptimized() error {
	couchdbURL := os.Getenv("COUCHDB_URL")
	if couchdbURL == "" {
		couchdbURL = "http://admin:password@localhost:5984/"
	}

	client, err := kivik.New("couch", couchdbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to CouchDB: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping CouchDB: %w", err)
	}
	OptimizedInstance = &OptimizedDB{
		Client:         client,
		maxConnections: 10, // Configurable connection pool size
		connPool:       make(chan *kivik.Client, 10),
		queryCache:     make(map[string]interface{}),
		queryMetrics:   make(map[string]*QueryMetrics),
	}

	// Initialize connection pool
	if err := OptimizedInstance.initConnectionPool(couchdbURL); err != nil {
		return fmt.Errorf("failed to initialize connection pool: %w", err)
	}

	// Create databases and indexes
	if err := OptimizedInstance.setupDatabases(); err != nil {
		return fmt.Errorf("failed to setup databases: %w", err)
	}

	// Create indexes for optimal query performance
	if err := OptimizedInstance.createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	utils.LogInfo("Optimized database initialized successfully", logrus.Fields{
		"connection_pool_size": OptimizedInstance.maxConnections,
		"databases":            []string{"posts", "users", "contacts"},
	})

	return nil
}

// initConnectionPool creates a pool of database connections
func (db *OptimizedDB) initConnectionPool(couchdbURL string) error {
	for i := 0; i < db.maxConnections; i++ {
		client, err := kivik.New("couch", couchdbURL)
		if err != nil {
			return fmt.Errorf("failed to create connection %d: %w", i, err)
		}
		db.connPool <- client
	}
	return nil
}

// getConnection gets a connection from the pool
func (db *OptimizedDB) getConnection() *kivik.Client {
	select {
	case conn := <-db.connPool:
		return conn
	default:
		// If pool is empty, return main client
		return db.Client
	}
}

// returnConnection returns a connection to the pool
func (db *OptimizedDB) returnConnection(conn *kivik.Client) {
	select {
	case db.connPool <- conn:
		// Connection returned to pool
	default:
		// Pool is full, connection will be garbage collected
	}
}

// setupDatabases creates databases if they don't exist
func (db *OptimizedDB) setupDatabases() error {
	ctx := context.Background()

	databases := []string{"posts", "users", "contacts"}

	for _, dbName := range databases {
		if exists, _ := db.Client.DBExists(ctx, dbName); !exists {
			if err := db.Client.CreateDB(ctx, dbName); err != nil {
				return fmt.Errorf("failed to create %s database: %w", dbName, err)
			}
			utils.LogInfo("Created database", logrus.Fields{"database": dbName})
		}
	}

	// Set database references
	db.PostsDB = db.Client.DB("posts")
	db.UsersDB = db.Client.DB("users")
	db.ContactsDB = db.Client.DB("contacts")

	return nil
}

// createIndexes creates optimized indexes for common queries
func (db *OptimizedDB) createIndexes() error {
	ctx := context.Background()

	// Posts indexes
	postsIndexes := []map[string]interface{}{
		{
			"index": map[string]interface{}{
				"fields": []string{"status", "published_at"},
			},
			"name": "status-published-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"author", "created_at"},
			},
			"name": "author-created-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"tags"},
			},
			"name": "tags-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"categories"},
			},
			"name": "categories-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"is_featured", "status"},
			},
			"name": "featured-status-index",
			"type": "json",
		},
	}

	for _, index := range postsIndexes {
		indexName := index["name"].(string)
		indexDef := index["index"]
		if err := db.PostsDB.CreateIndex(ctx, "", indexName, indexDef); err != nil {
			utils.LogError(err, "Failed to create posts index", logrus.Fields{
				"index": indexName,
			})
		} else {
			utils.LogInfo("Created posts index", logrus.Fields{
				"index": indexName,
			})
		}
	}

	// Users indexes
	usersIndexes := []map[string]interface{}{
		{
			"index": map[string]interface{}{
				"fields": []string{"username"},
			},
			"name": "username-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"email"},
			},
			"name": "email-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"role", "active"},
			},
			"name": "role-active-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"created_at"},
			},
			"name": "created-at-index",
			"type": "json",
		},
	}

	for _, index := range usersIndexes {
		indexName := index["name"].(string)
		indexDef := index["index"]
		if err := db.UsersDB.CreateIndex(ctx, "", indexName, indexDef); err != nil {
			utils.LogError(err, "Failed to create users index", logrus.Fields{
				"index": indexName,
			})
		} else {
			utils.LogInfo("Created users index", logrus.Fields{
				"index": indexName,
			})
		}
	}

	// Contacts indexes
	contactsIndexes := []map[string]interface{}{
		{
			"index": map[string]interface{}{
				"fields": []string{"status", "created_at"},
			},
			"name": "status-created-index",
			"type": "json",
		},
		{
			"index": map[string]interface{}{
				"fields": []string{"email"},
			},
			"name": "email-index",
			"type": "json",
		},
	}

	for _, index := range contactsIndexes {
		indexName := index["name"].(string)
		indexDef := index["index"]
		if err := db.ContactsDB.CreateIndex(ctx, "", indexName, indexDef); err != nil {
			utils.LogError(err, "Failed to create contacts index", logrus.Fields{
				"index": indexName,
			})
		} else {
			utils.LogInfo("Created contacts index", logrus.Fields{
				"index": indexName,
			})
		}
	}

	return nil
}

// trackQuery tracks query performance metrics
func (db *OptimizedDB) trackQuery(queryName string, duration time.Duration) {
	db.metricsMutex.Lock()
	defer db.metricsMutex.Unlock()

	if metrics, exists := db.queryMetrics[queryName]; exists {
		metrics.Count++
		metrics.TotalTime += duration
		metrics.AverageTime = metrics.TotalTime / time.Duration(metrics.Count)
		metrics.LastUsed = time.Now()
	} else {
		db.queryMetrics[queryName] = &QueryMetrics{
			Count:       1,
			TotalTime:   duration,
			AverageTime: duration,
			LastUsed:    time.Now(),
		}
	}
}

// GetQueryMetrics returns performance metrics for all queries
func (db *OptimizedDB) GetQueryMetrics() map[string]*QueryMetrics {
	db.metricsMutex.RLock()
	defer db.metricsMutex.RUnlock()

	// Create a copy to avoid race conditions
	metrics := make(map[string]*QueryMetrics)
	for k, v := range db.queryMetrics {
		metrics[k] = &QueryMetrics{
			Count:       v.Count,
			TotalTime:   v.TotalTime,
			AverageTime: v.AverageTime,
			LastUsed:    v.LastUsed,
		}
	}
	return metrics
}

// Health checks database connectivity and performance
func (db *OptimizedDB) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	if _, err := db.Client.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	duration := time.Since(start)
	if duration > 1*time.Second {
		utils.LogError(nil, "Database ping is slow", logrus.Fields{
			"duration": duration.String(),
		})
	}

	return nil
}

// Close closes all database connections
func (db *OptimizedDB) Close() error {
	// Close all connections in the pool
	close(db.connPool)
	for conn := range db.connPool {
		if err := conn.Close(); err != nil {
			utils.LogError(err, "Failed to close database connection", logrus.Fields{})
		}
	}

	// Close main client
	if err := db.Client.Close(); err != nil {
		return fmt.Errorf("failed to close main database client: %w", err)
	}

	utils.LogInfo("Database connections closed", logrus.Fields{})
	return nil
}
