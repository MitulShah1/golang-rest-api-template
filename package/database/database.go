package database

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	DriverName = "mysql"
	MaxDBConn  = 10
	MaxDBIdle  = 5
	DBTimeout  = 5 * time.Second
)

// DBConfig holds the database configuration
type DBConfig struct {
	Driver             string
	Host               string
	Port               string
	User               string
	Password           string
	DBName             string
	SSLMode            string
	MaxConn            int
	MaxIdle            int
	ConnectionTimeeout time.Duration
}

// Database wraps the sqlx.DB instance
type Database struct {
	DB *sqlx.DB
}

// NewDatabase initializes a new database connection
func NewDatabase(dbCnfg DBConfig) (*Database, error) {
	cfg := DBConfig{
		Host:     dbCnfg.Host,
		Port:     dbCnfg.Port,
		User:     dbCnfg.User,
		Password: dbCnfg.Password,
		DBName:   dbCnfg.DBName,
		SSLMode:  dbCnfg.SSLMode,
	}

	if cfg.Driver == "" {
		cfg.Driver = DriverName
	}

	dsn := getDSN(cfg)
	db, err := sqlx.Connect(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	if dbCnfg.MaxConn == 0 {
		dbCnfg.MaxConn = MaxDBConn
	}
	if dbCnfg.MaxIdle == 0 {
		dbCnfg.MaxIdle = MaxDBIdle
	}
	if dbCnfg.ConnectionTimeeout == 0 {
		dbCnfg.ConnectionTimeeout = DBTimeout
	}

	db.SetMaxOpenConns(dbCnfg.MaxConn)
	db.SetMaxIdleConns(dbCnfg.MaxIdle)
	db.SetConnMaxLifetime(dbCnfg.ConnectionTimeeout)

	return &Database{DB: db}, nil
}

// Close gracefully shuts down the database connection
func (d *Database) Close() {
	d.DB.Close()
}

// getDSN builds the Data Source Name based on the driver
func getDSN(cfg DBConfig) string {
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	case "sqlite3":
		return cfg.DBName // SQLite uses a simple file path as DSN
	default:
		log.Fatalf("Unsupported database driver: %s", cfg.Driver)
		return ""
	}
}
