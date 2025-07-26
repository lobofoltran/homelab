package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

type DBResult struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Status        string  `json:"status"`
	PingMS        float64 `json:"ping_ms,omitempty"`
	Version       string  `json:"version,omitempty"`
	DatabaseCount int     `json:"database_count,omitempty"`
	FoundDBA      bool    `json:"found_dba"`
	Error         string  `json:"error,omitempty"`
}

func CheckDatabases(central config.CentralConfig) {
	logger.Log.Info("[check_databases] Iniciando verificação de bancos de dados...")

	var results []DBResult

	for _, db := range config.Current.Databases {
		start := time.Now()
		dsn, driver := buildDSN(db)
		if dsn == "" {
			results = append(results, DBResult{Name: db.Name, Type: db.Type, Status: "unsupported"})
			continue
		}

		dbConn, err := sql.Open(driver, dsn)
		if err != nil {
			results = append(results, DBResult{Name: db.Name, Type: db.Type, Status: "failed", Error: err.Error()})
			continue
		}
		defer dbConn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err = dbConn.PingContext(ctx)
		if err != nil {
			results = append(results, DBResult{Name: db.Name, Type: db.Type, Status: "unreachable", Error: err.Error()})
			continue
		}

		pingMS := time.Since(start).Seconds() * 1000

		version := ""
		_ = dbConn.QueryRow("SELECT VERSION()").Scan(&version)

		foundDBA := false
		dbs := 0

		switch db.Type {
		case "mysql":
			foundDBA = existsDB(dbConn, "SHOW DATABASES", "_dba")
			dbs = countDatabases(dbConn, "SHOW DATABASES")
		case "postgres":
			foundDBA = existsDB(dbConn, "SELECT datname FROM pg_database", "_dba")
			dbs = countDatabases(dbConn, "SELECT datname FROM pg_database")
		case "sqlserver":
			foundDBA = existsDB(dbConn, "SELECT name FROM sys.databases", "_dba")
			dbs = countDatabases(dbConn, "SELECT name FROM sys.databases")
		}

		results = append(results, DBResult{
			Name:          db.Name,
			Type:          db.Type,
			Status:        "connected",
			PingMS:        pingMS,
			Version:       version,
			DatabaseCount: dbs,
			FoundDBA:      foundDBA,
		})
	}

	report := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"os":        runtime.GOOS,
		"results":   results,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("databases", report)
		logger.Log.Info("[check_databases] (dev) Resultado salvo localmente em /result")
	} else {
	}
}

func buildDSN(db config.DatabaseConfig) (string, string) {
	switch db.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.User, db.Password, db.Host, db.Port, db.Database), "mysql"
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, db.Host, db.Port, db.Database), "postgres"
	case "sqlserver":
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", db.User, db.Password, db.Host, db.Port, db.Database), "sqlserver"
	default:
		return "", ""
	}
}

func existsDB(db *sql.DB, query, target string) bool {
	rows, err := db.Query(query)
	if err != nil {
		return false
	}
	defer rows.Close()

	var name string
	for rows.Next() {
		_ = rows.Scan(&name)
		if name == target {
			return true
		}
	}
	return false
}

func countDatabases(db *sql.DB, query string) int {
	rows, err := db.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
	}
	return count
}
