package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/SwanHtetAungPhyo/grpcframe/pkg"
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Module management commands",
	Long:  "Commands for managing gRPC modules",
}

var moduleAddCmd = &cobra.Command{
	Use:   "add [module-name] [target-module]",
	Short: "Add a new gRPC module",
	Long:  "Creates a new gRPC module with handlers based on proto files",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		targetModule := args[1]
		if err := addModule(moduleName, targetModule); err != nil {
			pkg.Red.Printf("Failed to add module: %v\n", err)
			os.Exit(1)
		}
	},
}

type ModuleConfig struct {
	ModuleName     string
	TargetModule   string
	ProtoPath      string
	AppPath        string
	ProtogenPath   string
	ProtoFiles     []string
	ServiceMethods []ServiceMethod
}

type ServiceMethod struct {
	Name         string
	RpcName      string
	RequestType  string
	ResponseType string
	FileName     string
}

func addModule(moduleName, targetModule string) error {
	pkg.InfoLog("Starting to create new module", moduleName, "under target module", targetModule)

	config := &ModuleConfig{
		ModuleName:   moduleName,
		TargetModule: targetModule,
		ProtoPath:    filepath.Join("proto", moduleName),
		AppPath:      filepath.Join("app", "rpc", moduleName),
		ProtogenPath: "protogen",
	}

	pkg.InfoLog("Validating proto directory...")
	if err := validateProtoDirectory(config.ProtoPath); err != nil {
		return fmt.Errorf("proto validation failed: %w", err)
	}
	pkg.SuccessLog("Proto directory validated")

	pkg.InfoLog("Generating protobuf files...")
	if err := runProtogen(); err != nil {
		return fmt.Errorf("protogen failed: %w", err)
	}

	pkg.InfoLog("Discovering proto files...")
	if err := discoverProtoFiles(config); err != nil {
		return fmt.Errorf("failed to discover proto files: %w", err)
	}
	pkg.SuccessLog("Found", len(config.ProtoFiles), "proto files:", config.ProtoFiles)

	pkg.InfoLog("Extracting service methods...")
	if err := extractServiceMethods(config); err != nil {
		return fmt.Errorf("failed to extract service methods: %w", err)
	}
	pkg.SuccessLog("Extracted", len(config.ServiceMethods), "service methods")

	pkg.InfoLog("Creating module directory at", config.AppPath)
	if err := createModuleDirectory(config.AppPath); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	pkg.InfoLog("Generating service file...")
	if err := generateServiceFile(config); err != nil {
		return fmt.Errorf("failed to generate service file: %w", err)
	}

	pkg.InfoLog("Generating handler files...")
	if err := generateHandlerFiles(config); err != nil {
		return fmt.Errorf("failed to generate handler files: %w", err)
	}
	pkg.SuccessLog("Generated", len(config.ServiceMethods), "handler files")

	pkg.InfoLog("Generating converter file...")
	if err := generateConverterFile(config); err != nil {
		return fmt.Errorf("failed to generate converter file: %w", err)
	}

	pkg.InfoLog("Formatting generated code...")
	if err := exec.Command("go", "fmt", "./...").Run(); err != nil {
		return fmt.Errorf("failed to go fmt: %w", err)
	}
	pkg.InfoLog("Updating module dependencies...")
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return fmt.Errorf("failed to go mod tidy: %v", err.Error())
	}
	pkg.SuccessBox(fmt.Sprintf("Module '%s' created successfully!", moduleName))
	pkg.InfoLog("Next steps:")
	pkg.InfoLog("1. Implement business logic in the generated handlers")
	pkg.InfoLog("2. Register service in the gRPC server")
	pkg.InfoLog("3. Add necessary database models")

	return nil
}
func validateProtoDirectory(protoPath string) error {
	pkg.InfoLog(fmt.Sprintf("Checking proto directory at %s...", protoPath))
	if _, err := os.Stat(protoPath); os.IsNotExist(err) {
		return fmt.Errorf("proto directory does not exist: %s", protoPath)
	}

	files, err := filepath.Glob(filepath.Join(protoPath, "*.proto"))
	if err != nil {
		return fmt.Errorf("failed to list proto files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no proto files found in directory: %s", protoPath)
	}

	pkg.InfoLog(fmt.Sprintf("Found %d proto files", len(files)))
	return nil
}

func extractServiceMethods(config *ModuleConfig) error {
	pkg.InfoLog("Searching for generated gRPC files...")

	patterns := []string{
		filepath.Join(config.ProtogenPath, config.ModuleName, "*_grpc.pb.go"),
	}

	var grpcFiles []string
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		grpcFiles = append(grpcFiles, files...)
	}

	if len(grpcFiles) == 0 {
		return fmt.Errorf("no gRPC files found for module %s in %s", config.ModuleName, config.ProtogenPath)
	}

	pkg.InfoLog(fmt.Sprintf("Found %d gRPC definition files", len(grpcFiles)))

	for i, file := range grpcFiles {
		pkg.InfoLog(fmt.Sprintf("Parsing file %d/%d: %s", i+1, len(grpcFiles), filepath.Base(file)))
		methods, err := parseGrpcFile(file, config.ModuleName)
		if err != nil {
			pkg.WarningLog(fmt.Sprintf("Skipping file %s: %v", file, err))
			continue
		}
		config.ServiceMethods = append(config.ServiceMethods, methods...)
		pkg.InfoLog(fmt.Sprintf("Found %d methods in this file", len(methods)))
	}

	if len(config.ServiceMethods) == 0 {
		return fmt.Errorf("no service methods found for module %s", config.ModuleName)
	}

	pkg.SuccessLog(fmt.Sprintf("Total service methods discovered: %d", len(config.ServiceMethods)))
	return nil
}

func generateHandlerFiles(config *ModuleConfig) error {
	pkg.InfoLog(fmt.Sprintf("Generating %d handler files...", len(config.ServiceMethods)))

	for i, method := range config.ServiceMethods {
		pkg.InfoLog(fmt.Sprintf("Creating handler %d/%d: %s", i+1, len(config.ServiceMethods), method.Name))
		handlerContent := generateHandlerContent(config, method)
		handlerPath := filepath.Join(config.AppPath, method.FileName)
		if err := writeFile(handlerPath, handlerContent); err != nil {
			return fmt.Errorf("failed to write handler file %s: %w", handlerPath, err)
		}
	}
	return nil
}
func runProtogen() error {
	pkg.InfoLog("Running protogen...")
	cmd := exec.Command("make", "protoc")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make proto failed: %w", err)
	}
	pkg.InfoLog("Protogen completed successfully")
	return nil
}

func discoverProtoFiles(config *ModuleConfig) error {
	protoFiles, err := filepath.Glob(filepath.Join(config.ProtoPath, "*.proto"))
	if err != nil {
		return fmt.Errorf("failed to discover proto files: %w", err)
	}
	for _, file := range protoFiles {
		config.ProtoFiles = append(config.ProtoFiles, filepath.Base(file))
	}
	pkg.InfoLog(fmt.Sprintf("Discovered %d proto files", len(config.ProtoFiles)))
	return nil
}

// func extractServiceMethods(config *ModuleConfig) error {
//
//		patterns := []string{
//			filepath.Join(config.ProtogenPath, config.ModuleName, "*_grpc.pb.go"),
//		}
//
//		var grpcFiles []string
//		for _, pattern := range patterns {
//			files, err := filepath.Glob(pattern)
//			if err != nil {
//				continue
//			}
//			grpcFiles = append(grpcFiles, files...)
//		}
//
//		if len(grpcFiles) == 0 {
//			return fmt.Errorf("no gRPC files found for module %s in %s", config.ModuleName, config.ProtogenPath)
//		}
//
//		for _, file := range grpcFiles {
//			methods, err := parseGrpcFile(file, config.ModuleName)
//			if err != nil {
//				pkg.Red.Printf("Warning: failed to parse file %s: %v\n", file, err)
//				continue
//			}
//			config.ServiceMethods = append(config.ServiceMethods, methods...)
//		}
//
//		if len(config.ServiceMethods) == 0 {
//			return fmt.Errorf("no service methods found for module %s", config.ModuleName)
//		}
//
//		pkg.Info(fmt.Sprintf("Found %d service methods for module %s", len(config.ServiceMethods), config.ModuleName))
//		return nil
//	}
func parseGrpcFile(filePath, moduleName string) ([]ServiceMethod, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var methods []ServiceMethod
	lines := strings.Split(string(content), "\n")
	serviceInterfacePattern := regexp.MustCompile(`type\s+(\w+)Server\s+interface\s*{`)
	methodPattern := regexp.MustCompile(`^\s*(\w+)\s*\(\s*.*context\.Context\s*,\s*.*\*(\w+)\s*\)\s*\(\s*\*(\w+)\s*,\s*error\s*\)`)

	var inServiceInterface bool

	for _, line := range lines {
		if match := serviceInterfacePattern.FindStringSubmatch(line); match != nil {
			inServiceInterface = true
			continue
		}

		if inServiceInterface && strings.Contains(line, "}") {
			inServiceInterface = false
			continue
		}

		if inServiceInterface {
			if match := methodPattern.FindStringSubmatch(line); match != nil {
				methodName := match[1]
				requestType := match[2]
				responseType := match[3]
				fileName := fmt.Sprintf("rpc_%s.go", camelToSnake(methodName))

				methods = append(methods, ServiceMethod{
					Name:         methodName,
					RpcName:      methodName,
					RequestType:  requestType,
					ResponseType: responseType,
					FileName:     fileName,
				})
			}
		}
	}

	return methods, nil
}
func createModuleDirectory(appPath string) error {
	if err := os.MkdirAll(appPath, 0755); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}
	return nil
}

func generateServiceFile(config *ModuleConfig) error {
	serviceContent := generateServiceContent(config)
	servicePath := filepath.Join(config.AppPath, "service.go")
	if err := writeFile(servicePath, serviceContent); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}
	return nil
}

func generateServiceContent(config *ModuleConfig) string {
	packageName := strings.ToLower(config.ModuleName)
	serviceName := toPascalCase(config.ModuleName)
	pbImportPath := fmt.Sprintf("%s/protogen/%s", config.TargetModule, packageName)
	pbPackage := fmt.Sprintf("%spb", packageName)
	return fmt.Sprintf(`package %s

import (
	%s "%s"
)

type %sService struct {
	%s.Unimplemented%sServiceServer
}

func New%sService() *%sService {
	return &%sService{}
}
`, packageName, pbPackage, pbImportPath, serviceName, pbPackage, serviceName, serviceName, serviceName, serviceName)
}

func generateConverterFile(config *ModuleConfig) error {
	converterContent := generateConverterContent(config)
	converterPath := filepath.Join(config.AppPath, "converter.go")
	if err := writeFile(converterPath, converterContent); err != nil {
		return fmt.Errorf("failed to write converter file: %w", err)
	}
	return nil
}

func generateConverterContent(config *ModuleConfig) string {
	packageName := strings.ToLower(config.ModuleName)
	modelImportPath := fmt.Sprintf("%s/internal/repo", config.TargetModule)
	pbImportPath := fmt.Sprintf("%s/protogen/%s", config.TargetModule, packageName)

	return fmt.Sprintf(`package %s

import (
    "%s"
     %spb "%s"
 
)

// ProtoToModel converts protobuf message to database model
func ProtoToModel(pb *%spb.%s) *repo.%s {
    if pb == nil {
        return nil
    }
    
    return &repo.%s{
        // Add your conversion fields here
        // Example:
        // ID:        pb.GetId(),
        // Name:      pb.GetName(),
        // CreatedAt: time.Unix(pb.GetCreatedAt().GetSeconds(), int64(pb.GetCreatedAt().GetNanos())),
    }
}

// ModelToProto converts database model to protobuf message
func ModelToProto(model *repo.%s) *%spb.%s {
    if model == nil {
        return nil
    }
    
    return &%spb.%s{
        // Add your conversion fields here
        // Example:
        // Id:        model.ID,
        // Name:      model.Name,
        // CreatedAt: timestamppb.New(model.CreatedAt),
    }
}

// ModelsToProtos converts slice of models to slice of protos
func ModelsToProtos(models []*repo.%s) []*%spb.%s {
    protos := make([]*%spb.%s, len(models))
    for i, model := range models {
        protos[i] = ModelToProto(model)
    }
    return protos
}

// ProtosToModels converts slice of protos to slice of models
func ProtosToModels(protos []*%spb.%s) []*repo.%s {
    models := make([]*repo.%s, len(protos))
    for i, pb := range protos {
        models[i] = ProtoToModel(pb)
    }
    return models
}
`, packageName, modelImportPath, packageName, pbImportPath,
		packageName, toPascalCase(config.ModuleName), toPascalCase(config.ModuleName),
		toPascalCase(config.ModuleName),
		toPascalCase(config.ModuleName), packageName, toPascalCase(config.ModuleName),
		packageName, toPascalCase(config.ModuleName),
		toPascalCase(config.ModuleName), packageName, toPascalCase(config.ModuleName),
		packageName, toPascalCase(config.ModuleName),
		packageName, toPascalCase(config.ModuleName), toPascalCase(config.ModuleName),
		toPascalCase(config.ModuleName))
}

func generateHandlerContent(config *ModuleConfig, method ServiceMethod) string {
	packageName := strings.ToLower(config.ModuleName)
	serviceName := toPascalCase(config.ModuleName)
	pbImportPath := fmt.Sprintf("%s/protogen/%s", config.TargetModule, packageName)
	pbPackage := fmt.Sprintf("%spb", packageName)
	return fmt.Sprintf(`package %s

import (
	"context"
	%s "%s"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *%sService) %s(ctx context.Context, req *%s.%s) (*%s.%s, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &%s.%s{}
	return resp, nil
}
`, packageName, pbPackage, pbImportPath, serviceName, method.Name, pbPackage,
		method.RequestType, pbPackage, method.ResponseType, pbPackage, method.ResponseType)
}

func writeFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}
	return nil
}

func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || unicode.IsSpace(r)
	})
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(word[:1]))
			if len(word) > 1 {
				result.WriteString(strings.ToLower(word[1:]))
			}
		}
	}
	return result.String()
}

func init() {
	moduleCmd.AddCommand(moduleAddCmd)
	rootCmd.AddCommand(moduleCmd)
}
