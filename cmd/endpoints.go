package cmd

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/grpcframe/pkg"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var gatewayRegisterCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Register all gRPC gateway endpoints",
	Long:  "Automatically discovers and registers all gRPC gateway endpoints in the gateway.go file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := registerGatewayEndpoints(); err != nil {
			pkg.Red.Printf("Failed to register gateway endpoints: %v\n", err)
			os.Exit(1)
		}
	},
}

type GatewayRegistration struct {
	ModuleName   string
	ServiceName  string
	PbImportPath string
	PbPackage    string
	RegisterFunc string
}

func registerGatewayEndpoints() error {
	// Discover all modules in app/rpc directory
	modules, err := discoverGatewayModules()
	if err != nil {
		return fmt.Errorf("failed to discover modules: %w", err)
	}

	if len(modules) == 0 {
		pkg.InfoLog("No modules found to register in gateway")
		return nil
	}

	// Generate updated gateway.go content
	gatewayContent, err := generateGatewayContent(modules)
	if err != nil {
		return fmt.Errorf("failed to generate gateway content: %w", err)
	}

	// Write updated gateway.go
	gatewayPath := filepath.Join("app", "gateway", "gateway.go")
	if err := writeFile(gatewayPath, gatewayContent); err != nil {
		return fmt.Errorf("failed to write gateway file: %w", err)
	}

	// Format the file
	// TODO: Replace exec.Command with proper implementation
	// Original line: // TODO: Replace exec.Command with proper implementation
	// TODO: Replace exec.Command with proper implementation
	// Original line: // Original line: 	if err := exec.Command("gofmt", "-w", gatewayPath).Run(); err != nil {

	// Run go mod tidy
	// TODO: Replace exec.Command with proper implementation
	// Original line: // TODO: Replace exec.Command with proper implementation
	// TODO: Replace exec.Command with proper implementation
	// Original line: // Original line: 	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {

	pkg.InfoLog(fmt.Sprintf("Successfully registered %d gateway endpoints", len(modules)))
	return nil
}

func discoverGatewayModules() ([]GatewayRegistration, error) {
	var registrations []GatewayRegistration

	// Get target module name from go.mod
	targetModule, err := getTargetModuleName()
	if err != nil {
		return nil, fmt.Errorf("failed to get target module name: %w", err)
	}

	rpcPath := filepath.Join("app", "rpc")
	if _, err := os.Stat(rpcPath); os.IsNotExist(err) {
		return registrations, nil
	}

	// Read all directories in app/rpc (excluding server.go)
	entries, err := os.ReadDir(rpcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rpc directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip files like server.go
		}

		moduleName := entry.Name()

		// Check if service.go exists
		servicePath := filepath.Join(rpcPath, moduleName, "service.go")
		if _, err := os.Stat(servicePath); os.IsNotExist(err) {
			continue
		}

		// Extract gateway registration information
		gatewayInfo, err := extractGatewayInfo(servicePath, moduleName, targetModule)
		if err != nil {
			pkg.Red.Printf("Warning: failed to extract gateway info for %s: %v\n", moduleName, err)
			continue
		}

		registrations = append(registrations, gatewayInfo)
	}

	return registrations, nil
}

func extractGatewayInfo(servicePath, moduleName, targetModule string) (GatewayRegistration, error) {
	packageName := strings.ToLower(moduleName)
	serviceName := toPascalCase(moduleName)
	pbPackage := fmt.Sprintf("%spb", packageName)

	return GatewayRegistration{
		ModuleName:   moduleName,
		ServiceName:  serviceName,
		PbImportPath: fmt.Sprintf("%s/protogen/%s", targetModule, packageName),
		PbPackage:    pbPackage,
		RegisterFunc: fmt.Sprintf("%s.Register%sServiceHandlerFromEndpoint", pbPackage, serviceName),
	}, nil
}

func generateGatewayContent(registrations []GatewayRegistration) (string, error) {
	_, err := getTargetModuleName()
	if err != nil {
		return "", err
	}

	var content strings.Builder

	// Package declaration and imports
	content.WriteString(`package gateway

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
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
`)

	// Add proto imports
	for _, reg := range registrations {
		content.WriteString(fmt.Sprintf("    %s \"%s\"\n", reg.PbPackage, reg.PbImportPath))
	}

	content.WriteString(`    "github.com/rakyll/statik/fs"
    "github.com/rs/cors"
    "github.com/sirupsen/logrus"
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

    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(25 * 1024 * 1024)), // 25MB
    }

`)

	// Add service registrations
	for _, reg := range registrations {
		content.WriteString(fmt.Sprintf(`    if err := %s(ctx, gwMux, g.grpcAddr, opts); err != nil {
        return fmt.Errorf("failed to register %s service gateway: %%w", err)
    }

`, reg.RegisterFunc, strings.ToLower(reg.ServiceName)))
	}

	// Add the rest of the Start method
	content.WriteString(`    statikFS, err := fs.New()
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
`)

	return content.String(), nil
}

// Also update the server registration to use the correct path
func generateServerContent(registrations []ServiceRegistration) (string, error) {
	targetModule, err := getTargetModuleName()
	if err != nil {
		return "", err
	}

	var content strings.Builder

	content.WriteString(`package rpc

import (
    db "` + targetModule + `/internal/repo"
`)

	for _, reg := range registrations {
		content.WriteString(fmt.Sprintf("    %ssv \"%s/app/rpc/%s\"\n",
			strings.ToLower(reg.ServiceName), targetModule, reg.ModuleName))
		content.WriteString(fmt.Sprintf("    %s \"%s\"\n",
			reg.PbPackage, reg.PbImportPath))
	}

	content.WriteString(`    "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
    "net"
)

// Server implements the gRPC services
type Server struct {
`)

	for _, reg := range registrations {
		content.WriteString(fmt.Sprintf("    %s.Unimplemented%sServiceServer\n",
			reg.PbPackage, reg.ServiceName))
	}

	content.WriteString(`    store  *db.Store
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

// Run starts the gRPC server
func (s *Server) Run() error {
    listener, err := net.Listen("tcp", ":9001")
    if err != nil {
        panic(err.Error())
    }
    grpcServer := grpc.NewServer()

`)

	for _, reg := range registrations {
		content.WriteString(fmt.Sprintf("    %s := %ssv.New%sService()\n",
			reg.ServiceVar, strings.ToLower(reg.ServiceName), reg.ServiceName))
	}

	content.WriteString("\n    // Register services with gRPC server\n")

	for _, reg := range registrations {
		content.WriteString(fmt.Sprintf("    %s(grpcServer, %s)\n",
			reg.RegisterFunc, reg.ServiceVar))
	}

	content.WriteString(`
    s.logger.Println("Starting server on port 9001")
    err = grpcServer.Serve(listener)
    if err != nil {
        s.logger.Fatal(err.Error())
        return err
    }
    return nil
}`)

	return content.String(), nil
}

// Update the register services function to use correct path
func registerServices() error {
	// Discover all modules in app/rpc directory
	modules, err := discoverModules()
	if err != nil {
		return fmt.Errorf("failed to discover modules: %w", err)
	}

	if len(modules) == 0 {
		pkg.InfoLog("No modules found to register")
		return nil
	}

	// Generate updated server.go content
	serverContent, err := generateServerContent(modules)
	if err != nil {
		return fmt.Errorf("failed to generate server content: %w", err)
	}

	// Write updated server.go - note the correct path
	serverPath := filepath.Join("app", "rpc", "server.go")
	if err := writeFile(serverPath, serverContent); err != nil {
		return fmt.Errorf("failed to write server file: %w", err)
	}

	// Format the file
	// TODO: Replace exec.Command with proper implementation
	// Original line: // TODO: Replace exec.Command with proper implementation
	// TODO: Replace exec.Command with proper implementation
	// Original line: // Original line: 	if err := exec.Command("gofmt", "-w", serverPath).Run(); err != nil {

	// Run go mod tidy
	// TODO: Replace exec.Command with proper implementation
	// Original line: // TODO: Replace exec.Command with proper implementation
	// TODO: Replace exec.Command with proper implementation
	// Original line: // Original line: 	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {

	pkg.InfoLog(fmt.Sprintf("Successfully registered %d services", len(modules)))
	return nil
}

// Update init function
func init() {
	moduleCmd.AddCommand(moduleAddCmd)
	moduleCmd.AddCommand(moduleRegisterCmd)
	moduleCmd.AddCommand(gatewayRegisterCmd)
	rootCmd.AddCommand(moduleCmd)
}
