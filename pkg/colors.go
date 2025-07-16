// pkg/logger.go
package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	Red     = color.New(color.FgRed)
	Green   = color.New(color.FgGreen)
	Yellow  = color.New(color.FgYellow)
	Blue    = color.New(color.FgBlue)
	Cyan    = color.New(color.FgCyan)
	White   = color.New(color.FgWhite)
	Bold    = color.New(color.Bold)
	Success = color.New(color.Bold, color.FgGreen)
	Warning = color.New(color.Bold, color.FgYellow)
	Error   = color.New(color.Bold, color.FgRed)
	Info    = color.New(color.Bold, color.FgBlue)
)

// Print functions with error handling
func printColor(c *color.Color, a ...interface{}) {
	_, err := c.Println(a...)
	if err != nil {
		Red.Println("logging error:", err)
	}
}

func printfColor(c *color.Color, format string, a ...interface{}) {
	_, err := c.Printf(format, a...)
	if err != nil {
		Red.Println("logging error:", err)
	}
}

// Logging functions
func InfoLog(val ...interface{}) {
	printColor(Info, append([]interface{}{"[INFO]"}, val...)...)
}

func SuccessLog(val ...interface{}) {
	printColor(Success, append([]interface{}{"[SUCCESS]"}, val...)...)
}

func WarningLog(val ...interface{}) {
	printColor(Warning, append([]interface{}{"[WARNING]"}, val...)...)
}

func ErrorLog(val ...interface{}) {
	printColor(Error, append([]interface{}{"[ERROR]"}, val...)...)
}

// Boxed success message
func SuccessBox(message string) {
	line := "╔" + strings.Repeat("═", len(message)+4) + "╗"
	middle := "║  " + message + "  ║"
	bottom := "╚" + strings.Repeat("═", len(message)+4) + "╝"

	Success.Println(line)
	Success.Println(middle)
	Success.Println(bottom)
}

// Section header
func Section(title string) {
	fmt.Println()
	Info.Printf("== %s ==\n", strings.ToUpper(title))
	fmt.Println()
}

// Progress indicator
func Progress(current, total int, message string) {
	percent := float64(current) / float64(total) * 100
	Cyan.Printf("[%d/%d %.0f%%] %s\n", current, total, percent, message)
}

// Fatal error
func Fatal(err error) {
	ErrorLog(err)
	os.Exit(1)
}

// Debug logging (conditional)
var DebugEnabled = false

func DebugLog(val ...interface{}) {
	if DebugEnabled {
		printColor(Cyan, append([]interface{}{"[DEBUG]"}, val...)...)
	}
}
