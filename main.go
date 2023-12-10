package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/seigaalghi/seitest/generator"
	"github.com/seigaalghi/seitest/utils"
	"github.com/spf13/cobra"
)

var message string
var forceFlag bool
var rootCmd = &cobra.Command{
	Use:   "seitest",
	Short: "Sei's Unit Testing Generator",
	Long:  "Unit Testing Generator CLI created by Seiga",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Seitest Installed!")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Seitest",
	Run: func(cmd *cobra.Command, args []string) {
		command := exec.Command("go", "install", "github.com/vektra/mockery/v2@v2.38.0")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}

		command = exec.Command("go", "install", "golang.org/x/tools/cmd/goimports@latest")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}

		command = exec.Command("mockery", "--all")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}

	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Unit Test",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Println("file path is required")
			os.Exit(1)
		}

		forced, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Println("failed : ", err.Error())
		}

		functions, structs, imports, err := utils.ScanFunctions(args[0])
		if err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}

		var executed []string
		for _, f := range functions {
			path := strings.Replace(f.FilePath, ".go", "_test.go", 1)
			if !forced {
				if utils.FileExists(path) && !utils.InArray(executed, path) {
					return
				}
			}

			var lines []string
			var file *os.File
			var err error
			if !utils.InArray(executed, path) {
				file, err = os.Create(path)
				if err != nil {
					fmt.Println("Failed creating file", err.Error())
					os.Exit(1)
				}
				executed = append(executed, path)
				lines = append(lines, fmt.Sprintf("package %s", f.Package))
				lines = append(lines, `import "testing"`)
				for _, imp := range imports {
					if imp.FilePath == f.FilePath {
						lines = append(lines, imp.Content+"\n")
					}
				}
			} else {
				file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Failed opening file", err.Error())
					os.Exit(1)
				}
			}

			writer := bufio.NewWriter(file)

			switch f.IsMethod {
			case true:
				lines, _ = generator.MethodTestGenerator(f, lines, structs)
			default:
				lines, _ = generator.FuncTestGenerator(f, lines)
			}

			for _, line := range lines {
				writer.WriteString(line + "\n")
			}

			err = writer.Flush()
			if err != nil {
				log.Fatal(err.Error())
			}
			file.Close()
		}

		for _, exe := range executed {
			command := exec.Command("go", "fmt", exe)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Run(); err != nil {
				fmt.Println("Error executing command:", err)
				os.Exit(1)
			}
			command = exec.Command("goimports", "-w", exe)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Run(); err != nil {
				fmt.Println("Error executing command:", err)
				os.Exit(1)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(generateCmd)
	initCmd.Flags().StringVarP(&message, "message", "m", "Hello, World!", "Message to display")
	generateCmd.Flags().StringVarP(&message, "generate", "g", "", "Generate unit test")
	generateCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force overwrite")

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
