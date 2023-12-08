package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	"github.com/seigaalghi/seitest/utils"
	"github.com/spf13/cobra"
)

var message string
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
				FuncTestGenerator(f, &executed)
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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func InArray(list []string, file string) bool {
	for _, l := range list {
		if l == file {
			return true
		}
	}

	return false
}

func FuncTestGenerator(f utils.Function, executed *[]string) {
	path := strings.Replace(f.FilePath, ".go", "_test.go", 1)
	var lines []string
	var file *os.File
	var err error
	if !InArray(*executed, path) {
		file, err = os.Create(path)
		if err != nil {
			fmt.Println("Failed creating file", err.Error())
			os.Exit(1)
		}
		*executed = append(*executed, path)
		lines = append(lines, fmt.Sprintf("package %s", f.Package))
		lines = append(lines, `import "testing"`)
	} else {
		file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Failed opening file", err.Error())
			os.Exit(1)
		}
	}

	writer := bufio.NewWriter(file)

	write(lines, f, writer)

	err = writer.Flush()
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
}

func write(lines []string, f utils.Function, writer *bufio.Writer) error {
	_ = strings.Split(f.Content, "\n")
	lines = append(lines, "\n")
	lines = append(lines, fmt.Sprintf(`func Test%s(t *testing.T){`, f.Name))
	defer func() {
		lines = append(lines, "\n}")
		for _, line := range lines {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}()

	lines = append(lines, parsePayloadToStruct(f.Payload))
	lines = append(lines, parseResponseToAssert(f.Result))
	lines = append(lines, parseTestRunner(f.Payload, f.Result, f.Name))

	return nil
}

func parsePayloadToStruct(payload string) string {

	// Create the payload struct with the formatted output
	payloadStruct := fmt.Sprintf("type payload struct {\n%s}", payloadToNamedVariable(payload))

	return payloadStruct
}

func payloadToNamedVariable(payload string) string {
	pieces := strings.Split(removeSuffixPrefixParentheses(payload), ", ")

	occupied := map[string]string{}
	var outputText string
	for _, piece := range pieces {
		fields := strings.Fields(piece) // Split by space to get individual words
		var fieldType string
		if len(fields) != 2 {
			fieldType = fields[0]
		} else {
			fieldType = fields[1]
		}
		initial := abbreviateString(fieldType)
		final := initial
		for i := 1; occupied[initial+strconv.Itoa(i)] != ""; i++ {
			final = initial + strconv.Itoa(i)
		}
		outputText += fmt.Sprintf("%s %s\n", final, fieldType)
	}

	return outputText
}

func parseResponseToAssert(response string) string {
	responseArr := strings.Split(removeSuffixPrefixParentheses(response), ",")
	for i, r := range responseArr {
		if splitted := strings.Split(r, " "); len(splitted) == 2 {
			responseArr[i] = splitted[1]
		}
	}
	return fmt.Sprintf(`tests := []struct {
		scenario string
		payload  payload
		assert func(*testing.T, %s)
	}{}`, strings.Join(responseArr, ", "))
}

func parseTestRunner(payload, response, functionName string) string {
	responseArr := strings.Split(removeSuffixPrefixParentheses(response), ",")
	for i := range responseArr {
		responseArr[i] = intToExcelColumn(i)
	}

	res := strings.Join(responseArr, ", ")

	payloadArr := strings.Split(removeSuffixPrefixParentheses(payload), ", ")
	for i, r := range payloadArr {
		if splitted := strings.Split(r, " "); len(splitted) == 2 {
			p := fmt.Sprintf("tt.payload.%s", abbreviateString(splitted[1]))
			payloadArr[i] = p
		} else {
			p := fmt.Sprintf("tt.payload.%s", abbreviateString(splitted[0]))
			payloadArr[i] = p
		}
	}

	payloadArr = addNumbersToDuplicates(payloadArr)
	pay := strings.Join(payloadArr, ", ")

	return fmt.Sprintf(`for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			%s := %s(%s)
			tt.assert(t, %s)
		})
	}`, res, functionName, pay, res)
}

func intToExcelColumn(n int) string {
	if n < 0 {
		return ""
	}

	var result strings.Builder

	for n >= 0 {
		result.WriteByte(byte('a' + (n % 26)))
		if n < 25 {
			break
		}
		n = (n / 26) - 1
	}

	// Reverse the string
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func abbreviateString(input string) string {
	inputArr := strings.Split(input, ".")
	if len(inputArr) == 2 {
		input = inputArr[1]
	}

	var abbreviation strings.Builder
	abbreviation.WriteRune(unicode.ToLower(rune(input[0])))

	for i := 1; i < len(input); i++ {
		if unicode.IsUpper(rune(input[i])) {
			abbreviation.WriteRune(unicode.ToLower(rune(input[i])))
		}
	}

	return abbreviation.String()
}

func removeSuffixPrefixParentheses(input string) string {
	input = strings.TrimPrefix(input, "(")
	input = strings.TrimSuffix(input, ")")

	return input
}

func addNumbersToDuplicates(arr []string) []string {
	// Create a map to keep track of element occurrences
	countMap := make(map[string]int)

	result := make([]string, len(arr))

	for i, element := range arr {
		if countMap[element] == 0 {
			// If the element is unique, keep it as it is
			result[i] = element
		} else {
			// If the element is a duplicate, append the count to it
			result[i] = fmt.Sprintf("%s%d", element, countMap[element])
		}

		// Increment the count for the element
		countMap[element]++
	}

	return result
}