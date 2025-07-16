package cmd

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/grpcframe/pkg"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var moduleRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register all gRPC services in server.go",
	Long:  "Automatically discovers and registers all gRPC services in the server.go file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := registerServices(); err != nil {
			pkg.Red.Printf("Failed to register services: %v\n", err)
			os.Exit(1)
		}
	},
}

type ServiceRegistration struct {
	ModuleName   string
	ServiceName  string
	PbImportPath string
	PbPackage    string
	ServiceVar   string
	RegisterFunc string
}

func discoverModules() ([]ServiceRegistration, error) {
	var registrations []ServiceRegistration

	// Get target module name from go.mod
	targetModule, err := getTargetModuleName()
	if err != nil {
		return nil, fmt.Errorf("failed to get target module name: %w", err)
	}

	rpcPath := filepath.Join("app", "rpc")
	if _, err := os.Stat(rpcPath); os.IsNotExist(err) {
		return registrations, nil
	}

	// Read all directories in app/rpc
	entries, err := os.ReadDir(rpcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rpc directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		moduleName := entry.Name()

		// Check if service.go exists
		servicePath := filepath.Join(rpcPath, moduleName, "service.go")
		if _, err := os.Stat(servicePath); os.IsNotExist(err) {
			continue
		}

		// Extract service information
		serviceInfo, err := extractServiceInfo(servicePath, moduleName, targetModule)
		if err != nil {
			pkg.Red.Printf("Warning: failed to extract service info for %s: %v\n", moduleName, err)
			continue
		}

		registrations = append(registrations, serviceInfo)
	}

	return registrations, nil
}

func getTargetModuleName() (string, error) {
	goModPath := "go.mod"
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func extractServiceInfo(servicePath, moduleName, targetModule string) (ServiceRegistration, error) {
	_, err := os.ReadFile(servicePath)
	if err != nil {
		return ServiceRegistration{}, fmt.Errorf("failed to read service file: %w", err)
	}

	packageName := strings.ToLower(moduleName)
	serviceName := toPascalCase(moduleName)
	pbPackage := fmt.Sprintf("%spb", packageName)

	return ServiceRegistration{
		ModuleName:   moduleName,
		ServiceName:  serviceName,
		PbImportPath: fmt.Sprintf("%s/protogen/%s", targetModule, packageName),
		PbPackage:    pbPackage,
		ServiceVar:   fmt.Sprintf("%sService", strings.ToLower(serviceName)),
		RegisterFunc: fmt.Sprintf("%s.Register%sServiceServer", pbPackage, serviceName),
	}, nil
}

//func generateServerContent(registrations []ServiceRegistration) (string, error) {
//	targetModule, err := getTargetModuleName()
//	if err != nil {
//		return "", err
//	}
//
//	var imports strings.Builder
//	var serviceVars strings.Builder
//	var serviceInits strings.Builder
//	var serviceRegistrations strings.Builder
//
//	imports.WriteString(`package api
//
//import (
//    db "` + targetModule + `/internal/repo"
//`)
//
//	for _, reg := range registrations {
//		imports.WriteString(fmt.Sprintf("    %ssv \"%s/app/rpc/%s\"\n",
//			strings.ToLower(reg.ServiceName), targetModule, reg.ModuleName))
//		imports.WriteString(fmt.Sprintf("    %s \"%s\"\n",
//			reg.PbPackage, reg.PbImportPath))
//	}
//
//	imports.WriteString(`    "github.com/sirupsen/logrus"
//    "google.golang.org/grpc"
//    "net"
//)
//
//// Server implements the gRPC services
//type Server struct {
//`)
//
//	for _, reg := range registrations {
//		serviceVars.WriteString(fmt.Sprintf("    %s.Unimplemented%sServiceServer\n",
//			reg.PbPackage, reg.ServiceName))
//	}
//
//	serviceVars.WriteString(`    store  *db.Store
//    logger *logrus.Logger
//}
//
//func NewServer(
//    db *db.Store,
//    logger *logrus.Logger,
//) *Server {
//    return &Server{
//        store:  db,
//        logger: logger,
//    }
//}
//
//// Run starts the gRPC server
//func (s *Server) Run() error {
//    listener, err := net.Listen("tcp", ":9001")
//    if err != nil {
//        panic(err.Error())
//    }
//    grpcServer := grpc.NewServer()
//
//`)
//
//	for _, reg := range registrations {
//		serviceInits.WriteString(fmt.Sprintf("    %s := %ssv.New%sService()\n",
//			reg.ServiceVar, strings.ToLower(reg.ServiceName), reg.ServiceName))
//	}
//
//	serviceInits.WriteString("\n    // Register services with gRPC server\n")
//
//	for _, reg := range registrations {
//		serviceRegistrations.WriteString(fmt.Sprintf("    %s(grpcServer, %s)\n",
//			reg.RegisterFunc, reg.ServiceVar))
//	}
//
//	serviceRegistrations.WriteString(`
//    s.logger.Println("Starting server on port 9001")
//    err = grpcServer.Serve(listener)
//    if err != nil {
//        s.logger.Fatal(err.Error())
//        return err
//    }
//    return nil
//}`)
//
//	return imports.String() + serviceVars.String() + serviceInits.String() + serviceRegistrations.String(), nil
//}

func init() {
	moduleCmd.AddCommand(moduleAddCmd)
	moduleCmd.AddCommand(moduleRegisterCmd)
	rootCmd.AddCommand(moduleCmd)
}
