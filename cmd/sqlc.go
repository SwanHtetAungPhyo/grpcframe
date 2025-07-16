package cmd

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/SwanHtetAungPhyo/grpcframe/pkg"
	"github.com/spf13/cobra"
)

var sqlcCmd = &cobra.Command{
	Use:   "sqlc",
	Short: "SQLc generation commands",
	Long:  "Commands for running SQLc code generation",
}

var sqlcGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Run SQLc code generation for all queries",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runSQLcGenerate(); err != nil {
			pkg.Red.Printf("Failed to generate SQLc code: %v\n", err)
		} else {
			pkg.InfoLog("SQLc code generated successfully")
		}
	},
}

func runSQLcGenerate() error {
	cmd := exec.Command("sqlc", "generate")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlc generate failed: %w\n%s", err, stderr.String())
	}
	err := exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}
	pkg.InfoLog(stdout.String())
	return nil
}

func init() {
	sqlcCmd.AddCommand(sqlcGenerateCmd)
	rootCmd.AddCommand(sqlcCmd)
}
