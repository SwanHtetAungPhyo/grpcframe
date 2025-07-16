package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/SwanHtetAungPhyo/grpcframe/pkg"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [dir] [module-name]",
	Short: "Initialize a new project",
	Long:  "Initializes a new project at the specified path",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		moduleName := args[1]

		if err := initializeProject(path, moduleName); err != nil {
			pkg.Red.Printf("Failed to initialize project: %v\n", err)
			os.Exit(1)
		}
	},
}

// ProjectConfig holds configuration for project initialization
type ProjectConfig struct {
	ProjectPath string
	ModuleName  string
	GoVersion   string
}

// initializeProject orchestrates the entire project initialization process
func initializeProject(projectPath, moduleName string) error {
	config := &ProjectConfig{
		ProjectPath: projectPath,
		ModuleName:  moduleName,
		GoVersion:   getGoVersion(),
	}

	// Create project directory if it doesn't exist
	if err := createProjectDirectory(config.ProjectPath); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create directory structure
	if err := createDirectoryStructure(config.ProjectPath); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Initialize Go module
	if err := initializeGoModule(config.ProjectPath, config.ModuleName); err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}

	// Create project files
	if err := createProjectFiles(config); err != nil {
		return fmt.Errorf("failed to create project files: %w", err)
	}

	//pkg.InfoLog("Go mod tidying completed successfully")
	//if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
	//	return fmt.Errorf("failed to go mod tidy: %w", err)
	//}
	pkg.InfoLog(" Please  manually run go mod tidy")
	pkg.SuccessBox(fmt.Sprintf("Module '%s' created successfully!", moduleName))

	pkg.InfoLog("Project initialization completed successfully")
	return nil
}

// createProjectDirectory creates the main project directory
func createProjectDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		pkg.InfoLog("Creating " + path)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", path, err)
		}
	}
	return nil
}

// createDirectoryStructure creates all required subdirectories
func createDirectoryStructure(baseDir string) error {
	dirs := []string{
		"proto",
		"protogen",
		"app",
		"app/rpc",
		"app/gateway",
		"cmd",
		"internal",
		"internal/repo",
		"internal/services",
		"database",
		"database/migration",
		"database/schema",
		"database/queries",
		"pkg",
		"pkg/utils",
		"doc/swagger",
		"pkg/utils/convert",
		"pkg/utils/env",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(baseDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			pkg.InfoLog("Creating " + fullPath)
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %w", fullPath, err)
			}
		}
	}
	return nil
}

// initializeGoModule initializes the Go module in the project directory
func initializeGoModule(projectPath, moduleName string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory %s: %w", projectPath, err)
	}

	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			pkg.Red.Printf("Warning: failed to return to original directory %s: %v\n", currentDir, err)
		}
	}()

	cmd := exec.Command("go", "mod", "init", moduleName)
	if err := cmd.Run(); err != nil {
		pkg.Red.Printf("Warning: failed to run go mod init: %s\n", err)
	}
	pkg.InfoLog("Go module " + moduleName + " initialized successfully")
	return nil
}

func getGoVersion() string {
	const defaultVersion = "1.21"

	output, err := exec.Command("go", "version").Output()
	if err != nil {
		return defaultVersion
	}
	re := regexp.MustCompile(`go(\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(output))

	if len(matches) > 1 {
		return matches[1]
	}

	pkg.Red.Printf("Warning: failed to parse Go version, using default %s\n", defaultVersion)
	return defaultVersion
}

func createProjectFiles(config *ProjectConfig) error {
	files := getFileTemplates(config)

	for filePath, content := range files {
		fullPath := filepath.Join(config.ProjectPath, filePath)

		if err := createFileWithContent(fullPath, content); err != nil {
			return fmt.Errorf("failed to create file %s: %w", fullPath, err)
		}
	}
	return nil
}

func createFileWithContent(filePath, content string) error {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		pkg.InfoLog("File already exists: " + filePath)
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			pkg.Red.Printf("Warning: failed to close file: %s\n", err)
		}
	}(file)

	if content != "" {
		if _, err := file.WriteString(content); err != nil {
			return fmt.Errorf("failed to write content: %w", err)
		}
	}

	pkg.InfoLog("Created file: " + filePath)
	return nil
}

func getFileTemplates(config *ProjectConfig) map[string]string {
	return map[string]string{
		"app/app.go":                     generateAppTemplates(config.ModuleName),
		"cmd/main.go":                    generateMainTemplate(config.ModuleName),
		"app/gateway/gateway.go":         generateGatewayTemplate(),
		"app/rpc/server.go":              getServerTemplate(config.ModuleName),
		"Dockerfile":                     generateDockerfile(config.GoVersion),
		"README.md":                      generateReadme(config.ModuleName, config.GoVersion),
		"sqlc.yaml":                      generateSqlcConfig(),
		"makefile":                       generateMakefile(),
		"buf.yaml":                       generateBufConfig(),
		".env":                           generateEnvFile(),
		"internal/repo/store.go":         generateStoreStruct(),
		"pkg/utils/convert/convertor.go": generateCommonConvertor(),
		"pkg/utils/env/envs.go":          generateEnvUtils(),
	}
}
func generateStoreStruct() string {
	return `
package repo

import "github.com/jackc/pgx/v5/pgxpool"

type Store struct {
	*Queries
	conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *Store {
	return &Store{
		conn: conn,
	}
}`
}
func generateEnvFile() string {
	return `# Add Your env variable`
}
func generateMainTemplate(moduleName string) string {
	return fmt.Sprintf(`
package main
	import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"%s/app"
	"%s/app/rpc"
	"%s/app/gateway"
	"%s/internal/repo"
	"%s/pkg/utils/env"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	grpcServerAddress := env.GetEnv("GRPC_SERVER_ADDRESS", ":9001")
	gprcGatewayAddress := env.GetEnv("GRPC_GATEWAY_ADDRESS", ":8082")
	dbConn := DatabaseConn(logger)
	dbStore := db.NewStore(dbConn)
	grpcServer := rpc.NewServer(dbStore, logger)
	grpcGateway := gateway.NewGateway(logger, grpcServerAddress, gprcGatewayAddress)
	server := app.NewApp(grpcServer, grpcGateway)
	err := server.Run()
	if err != nil {
		logger.WithError(err).Fatal("failed to start server")
		return
	}
}
func DatabaseConn(logger *logrus.Logger) *pgxpool.Pool {
	user := env.GetEnv("DB_USER", "postgres")
	password := env.GetEnv("DB_PASSWORD", "postgres")
	dbName := env.GetEnv("DB_NAME", "postgres")
	host := env.GetEnv("DB_HOST", "localhost")
	port := env.GetEnv("DB_PORT", "5432")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	)

	connPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		logger.Fatal(err)
	}
	return connPool
}

`, moduleName, moduleName, moduleName, moduleName, moduleName)
}
func generateAppTemplates(moduleName string) string {
	return fmt.Sprintf(`package app

import (
	api "%s/app/rpc"
	rpcGate "%s/app/rpc"
	"github.com/sirupsen/logrus"
	"sync"
)

type App struct {
	grpcServer  *api.Server
	grpcGateway *rpcGate.Gateway
	logger      *logrus.Logger
}

func NewApp(
	grpcServer *api.Server,
	grpcGateway *rpcGate.Gateway,
) *App {
	return &App{
		grpcServer:  grpcServer,
		grpcGateway: grpcGateway,
	}
}

func (app *App) Run() error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := app.grpcServer.Run(); err != nil {
			if app.logger != nil {
				app.logger.Error(err.Error())
			}
		}
	}()

	go func() {
		defer wg.Done()
		if err := app.grpcGateway.Start(); err != nil {
			if app.logger != nil {
				app.logger.Error(err.Error())
			}
		}
	}()

	wg.Wait()
	if app.logger != nil {
		app.logger.Info("server stopped")
	}
	return nil
}
`, moduleName, moduleName)
}
func generateGatewayTemplate() string {
	return `
package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	logger     *logrus.Logger
	grpcAddr   string
	httpAddr   string
	swaggerDir string
}

func NewGateway(logger *logrus.Logger, grpcAddr, httpAddr string) *Gateway {
	return &Gateway{
		logger:     logger,
		grpcAddr:   grpcAddr,
		httpAddr:   httpAddr,
		swaggerDir: "../doc/swagger",
	}
}

func (g *Gateway) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux(
		runtime.WithErrorHandler(g.errorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(g.headerMatcher),
	)

	//opts := []grpc.DialOption{
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(25 * 1024 * 1024)), // 25MB
	//}
	//
	

	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("statik filesystem error: %w", err)
	}

	// Create main mux router
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))
	mux.HandleFunc("/healthz", g.healthCheck)

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	// Configure HTTP server
	server := &http.Server{
		Addr:         g.httpAddr,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		g.logger.Info("Shutting down HTTP gateway...")
		if err := server.Shutdown(ctx); err != nil {
			g.logger.WithError(err).Error("HTTP gateway shutdown error")
		}
	}()

	g.logger.Infof("Starting HTTP gateway on %s (gRPC backend: %s)", g.httpAddr, g.grpcAddr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
		return fmt.Errorf("HTTP gateway start error: %w", err)
	}

	return nil
}

func (g *Gateway) errorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	g.logger.WithError(err).Error("gateway error")
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

func (g *Gateway) headerMatcher(key string) (string, bool) {
	switch key {
	case "X-Request-ID", "X-Correlation-ID":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func (g *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

`
}

func generateEnvUtils() string {
	return `
package env

import (
	"os"
	"strconv"
	"time"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
`
}
func getServerTemplate(moduleName string) string {
	return fmt.Sprintf(`
package rpc
import (	
	db "%s/internal/repo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	)
type Server struct {
	store  *db.Store  
	logger *logrus.Logger 
}

func NewServer(
	db *db.Store,
	logger *logrus.Logger,
) *Server {

	return &Server{
		store:  db,
		logger: logger,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer()

	s.logger.Println("Starting server on port 9001")
	err = grpcServer.Serve(listener)
	if err != nil {
		s.logger.Fatal(err.Error())
		return err
	}
	return nil
}


`, moduleName)
}

// generateDockerfile creates the Dockerfile content
func generateDockerfile(goVersion string) string {
	return fmt.Sprintf(`FROM golang:%s-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /app/server /server

EXPOSE 8082
EXPOSE 8083

CMD ["/server"]`, goVersion)
}

func generateCommonConvertor() string {
	return `package convertor


import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertStringToUUID(uuidStr string) (pgtype.UUID, error) {
	if uuidStr == "" {
		return pgtype.UUID{Valid: false}, nil
	}

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return pgtype.UUID{Valid: false}, err
	}

	var pgUUID pgtype.UUID
	pgUUID.Bytes = parsedUUID
	pgUUID.Valid = true

	return pgUUID, nil
}

// ConvertTimestamp converts pgtype.Timestamptz to *timestamppb.Timestamp
func ConvertTimestamp(ts pgtype.Timestamptz) *timestamppb.Timestamp {
	if !ts.Valid {
		return nil
	}
	return timestamppb.New(ts.Time)
}

// ConvertUUIDToString converts pgtype.UUID to string
func ConvertUUIDToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}

	u := uuid.UUID(pgUUID.Bytes)
	return u.String()
}

`
}

// generateReadme creates the README.md content
func generateReadme(moduleName, goVersion string) string {
	return fmt.Sprintf(`# %s

A gRPC-based microservice built with Go, featuring both gRPC API and HTTP gateway support.

## Project Structure

`+"```"+`
.
├── Dockerfile
├── README.md
├── app
│   ├── gateway
│   │   └── gateway.go
│   └── rpc-api
│       └── server.go
├── buf.yaml
├── cmd
│   └── main.go
├── database
│   ├── migration
│   ├── queries
│   └── schema
├── go.mod
├── internal
│   ├── repo
│   └── services
├── makefile
├── pkg
│   └── utils
├── proto
├── protogen
└── sqlc.yaml
`+"```"+`

## Directory Structure

- **app/**: Application layer containing API implementations
  - **gateway/**: HTTP gateway server implementation
  - **rpc-api/**: gRPC server implementation
- **cmd/**: Application entry point
- **database/**: Database-related files
  - **migration/**: Database migration files
  - **queries/**: SQL query files for sqlc
  - **schema/**: Database schema definitions
- **internal/**: Private application code
  - **repo/**: Repository layer for data access
  - **services/**: Business logic layer
- **pkg/**: Public library code and utilities
- **proto/**: Protocol buffer definitions
- **protogen/**: Generated protobuf code
- **buf.yaml**: Buf configuration for protobuf generation
- **sqlc.yaml**: SQLC configuration for database code generation

## Getting Started

### Prerequisites

- Go %s or later
- Docker (optional)
- Buf CLI (for protobuf generation)
- SQLC (for database code generation)

### Installation

1. Clone the repository
2. Install dependencies:
   `+"```bash"+`
   make deps
   `+"```"+`

### Development

#### Generate Code

Generate protobuf files:
`+"```bash"+`
make proto
`+"```"+`

Generate database code:
`+"```bash"+`
make sqlc
`+"```"+`

#### Build and Run

Build the application:
`+"```bash"+`
make build
`+"```"+`

Run the application:
`+"```bash"+`
make run
`+"```"+`

#### Testing

Run tests:
`+"```bash"+`
make test
`+"```"+`

Format code:
`+"```bash"+`
make fmt
`+"```"+`

Lint code:
`+"```bash"+`
make lint
`+"```"+`

### Docker

Build Docker image:
`+"```bash"+`
make docker-build
`+"```"+`

Run Docker container:
`+"```bash"+`
make docker-run
`+"```"+`

### Available Make Targets

Run `+"```make help```"+` to see all available targets.

## API Endpoints

- gRPC Server: :8082
- HTTP Gateway: :8083

## Configuration

Configuration can be set through environment variables or configuration files.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the MIT License.
`, moduleName, goVersion)
}

// generateSqlcConfig creates the sqlc.yaml content
func generateSqlcConfig() string {
	return `version: "2"
sql:
  - engine: "postgresql"
    queries: "./database/queries"
    schema: "./database/schema"
    gen:
      go:
        package: "db"
        out: "./internal/repo"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true`
}

// generateMakefile creates the Makefile content
func generateMakefile() string {
	return `# Variables
PROTO_DIR = proto
PROTOGEN_DIR = protogen
APP_DIR = app
BINARY_NAME = server
MAIN_FILE = cmd/main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

# Run the application
.PHONY: run
run:
	go run $(MAIN_FILE)

protoc: ## Generate protobuf files
	@echo "Generating proto..."
	cd proto && protoc --go_out=../protogen --go_opt=paths=source_relative \
	--go-grpc_out=../protogen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=../protogen --grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out=../doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=api \
	--experimental_allow_proto3_optional \
	./**/*.proto
	@echo "Generating swagger statik..."
	statik -src=doc/swagger -dest=doc -f
# Generate protobuf files
.PHONY: proto
proto:
	buf generate

# Generate database code with sqlc
.PHONY: sqlc
sqlc:
	sqlc generate

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	rm -rf $(PROTOGEN_DIR)/*

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Test
.PHONY: test
test:
	go test ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME) .

# Docker run
.PHONY: docker-run
docker-run:
	docker run -p 8082:8082 -p 8083:8083 $(BINARY_NAME)

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  proto       - Generate protobuf files"
	@echo "  sqlc        - Generate database code"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  help        - Show this help message"`
}

// generateBufConfig creates the buf.yaml content
func generateBufConfig() string {
	return `version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
build:
  roots:
    - proto
generate:
  - name: go
    out: protogen
    opt:
      - paths=source_relative
  - name: go-grpc
    out: protogen
    opt:
      - paths=source_relative
  - name: grpc-gateway
    out: protogen
    opt:
      - paths=source_relative
  - name: openapiv2
    out: protogen
    opt:
      - allow_merge=true
      - merge_file_name=api`
}

func init() {
	rootCmd.AddCommand(initCmd)
}
