package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"postgres"`
	Tables struct {
		Source      string `yaml:"source"`
		Destination string `yaml:"destination"`
		DateColumn  string `yaml:"date_column"`
		DateFormat  string `yaml:"date_format"`
	} `yaml:"tables"`
}

func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to SQLite
	sqliteDB, err := sql.Open("sqlite3", "./source.db")
	if err != nil {
		log.Fatalf("SQLite connection error: %v", err)
	}
	defer sqliteDB.Close()

	// Connect to PostgreSQL
	pg := config.Postgres
	pgConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode)

	pgDB, err := sql.Open("postgres", pgConn)
	if err != nil {
		log.Fatalf("PostgreSQL connection error: %v", err)
	}
	defer pgDB.Close()

	// Ensure destination table exists
	createTable := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		name TEXT,
		%s TIMESTAMP
	);`, config.Tables.Destination, config.Tables.DateColumn)

	if _, err := pgDB.Exec(createTable); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Query data from SQLite
	query := fmt.Sprintf("SELECT name, %s FROM %s", config.Tables.DateColumn, config.Tables.Source)
	rows, err := sqliteDB.Query(query)
	if err != nil {
		log.Fatalf("Failed to query SQLite: %v", err)
	}
	defer rows.Close()

	insertQuery := fmt.Sprintf(`INSERT INTO %s (name, %s) VALUES ($1, $2)`, config.Tables.Destination, config.Tables.DateColumn)

	for rows.Next() {
		var name, datetimeStr string
		if err := rows.Scan(&name, &datetimeStr); err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		parsedTime, err := time.Parse(config.Tables.DateFormat, datetimeStr)
		if err != nil {
			log.Printf("Time parse error: %v", err)
			continue
		}

		if _, err := pgDB.Exec(insertQuery, name, parsedTime); err != nil {
			log.Printf("Insert error: %v", err)
		} else {
			fmt.Printf("Migrated: %s at %s\n", name, parsedTime.Format(time.RFC3339))
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
