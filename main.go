package main

import (
	"fmt"
	"os"
	"os/exec"

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

		functions, err := utils.ScanFunctions(args[0])
		if err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}

		var executed []string

		for _, f := range functions {
			switch f.IsMethod {
			case true:
				continue
			default:
				generator.FuncTestGenerator(f, &executed, forced)
			}
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
