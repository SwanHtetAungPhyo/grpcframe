package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SwanHtetAungPhyo/grpcframe/pkg"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  "Commands for running database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrations("up", ""); err != nil {
			pkg.ErrorLog("Migration failed:", err)
			os.Exit(1)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down [number]",
	Short: "Rollback migrations (default: 1)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		steps := "1" // default to 1 step
		if len(args) > 0 {
			steps = args[0]
		}
		if err := runMigrations("down", steps); err != nil {
			pkg.ErrorLog("Rollback failed:", err)
			os.Exit(1)
		}
	},
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func loadDBConfig() (*DBConfig, error) {
	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file not found")
	}

	file, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &DBConfig{
		SSLMode: "disable", // default
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "DB_HOST":
			config.Host = value
		case "DB_PORT":
			config.Port = value
		case "DB_USER":
			config.User = value
		case "DB_PASSWORD":
			config.Password = value
		case "DB_NAME":
			config.Name = value
		case "DB_SSL_MODE":
			config.SSLMode = value
		}
	}

	if config.Host == "" || config.Port == "" || config.User == "" || config.Name == "" {
		return nil, fmt.Errorf("missing required database configuration in .env")
	}

	return config, nil
}

func createMigrateInstance() (*migrate.Migrate, error) {
	config, err := loadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load DB config: %w", err)
	}

	migrationsDir := filepath.Join("database", "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory not found at %s", migrationsDir)
	}

	// First verify the database connection
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
	)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	// Create migrate instance
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return m, nil
}

func runMigrations(direction string, steps string) error {
	m, err := createMigrateInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	switch direction {
	case "up":
		pkg.InfoLog("Running migrations up...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrations up: %w", err)
		}
		pkg.SuccessLog("Migrations up completed successfully")

	case "down":
		if steps == "" {
			steps = "1"
		}
		pkg.InfoLog(fmt.Sprintf("Rolling back %s migration(s)...", steps))
		if err := m.Steps(-1 * parseInt(steps)); err != nil {
			return fmt.Errorf("failed to run migrations down: %w", err)
		}
		pkg.SuccessLog("Migrations down completed successfully")

	default:
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	return nil
}

func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

var migrateForceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force database to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		if err := forceMigration(version); err != nil {
			pkg.ErrorLog("Force migration failed:", err)
			os.Exit(1)
		}
	},
}

func forceMigration(version string) error {
	m, err := createMigrateInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	pkg.InfoLog(fmt.Sprintf("Forcing database to version %s...", version))
	if err := m.Force(parseInt(version)); err != nil {
		return fmt.Errorf("failed to force migration: %w", err)
	}

	pkg.SuccessLog("Database version forced successfully")
	return nil
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		if err := showVersion(); err != nil {
			pkg.ErrorLog("Failed to get version:", err)
			os.Exit(1)
		}
	},
}

func showVersion() error {
	m, err := createMigrateInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	if dirty {
		pkg.WarningLog(fmt.Sprintf("Current version: %d (dirty)", version))
	} else {
		pkg.InfoLog(fmt.Sprintf("Current version: %d", version))
	}

	return nil
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateForceCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	rootCmd.AddCommand(migrateCmd)
}
