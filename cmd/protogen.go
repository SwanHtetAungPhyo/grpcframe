package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/SwanHtetAungPhyo/grpcframe/pkg"
	"github.com/spf13/cobra"
)

var protogenCmd = &cobra.Command{
	Use:   "protogen",
	Short: "Generate protobuf files",
	Long:  "Generates Go code from protobuf definitions using buf",
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeProtogen(); err != nil {
			pkg.Red.Printf("Protogen failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func executeProtogen() error {
	pkg.InfoLog("Starting protobuf code generation...")

	cmd := exec.Command("make", "protoc")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make proto command failed: %w", err)
	}

	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	pkg.InfoLog("Protobuf code generation completed successfully")
	return nil
}

func init() {
	rootCmd.AddCommand(protogenCmd)
}
