package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func NewConfig() Config {
	return Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_NAME"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
}

func NewConnection(cfg Config) (*sql.DB, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("host=%s port=%s ", cfg.Host, cfg.Port))
	builder.WriteString(fmt.Sprintf("user=%s password=%s ", cfg.User, cfg.Password))
	builder.WriteString(fmt.Sprintf("dbname=%s ", cfg.DBName))
	builder.WriteString("sslmode=disable ")

	params := builder.String()

	var err error
	var db *sql.DB
	counter := 0
	for db == nil && counter < 2 {
		db, err = sql.Open("postgres", params)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	runMigrations(cfg)

	return db, nil
}

func runMigrations(cfg Config) {
	dbUri := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	migrationsDir := "file://migrations"

	m, err := migrate.New(migrationsDir, dbUri)
	if err != nil {
		log.Print(err)
		return
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Print(err)
		return
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("Applied migrations: %d, Diry: %t\n", version, dirty)
}
