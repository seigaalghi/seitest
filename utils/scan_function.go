package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Function struct {
	FilePath string
	Name     string
	Content  string
	IsMethod bool
	Body     string
	Package  string
	Payload  string
	Result   string
}

func ScanFunctions(packagePath string) ([]Function, error) {
	var functions []Function

	// Create a token file set
	fset := token.NewFileSet()

	// Walk through the directory recursively
	err := filepath.Walk(packagePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			// Parse each Go file
			node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return err
			}

			// Extract function names
			for _, decl := range node.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok && !strings.HasSuffix(path, "_test.go") {
					funcText, err := getFileText(path, fn.Pos(), fn.End(), fset)
					if err != nil {
						log.Fatal(err.Error())
					}

					payload, err := getFileText(path, fn.Type.Params.Pos(), fn.Type.Params.End(), fset)
					if err != nil {
						log.Fatal(err.Error())
					}
					result, err := getFileText(path, fn.Type.Results.Pos(), fn.Type.Results.End(), fset)
					if err != nil {
						log.Fatal(err.Error())
					}

					functions = append(functions, Function{
						FilePath: path,
						Name:     fn.Name.Name,
						Content:  funcText,
						IsMethod: fn.Recv != nil,
						Package:  node.Name.Name,
						Payload:  payload,
						Result:   result,
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return functions, nil
}

func getFileText(filePath string, start, end token.Pos, fset *token.FileSet) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	fileText := string(data)

	// Convert token.Pos offsets to file positions
	startPos := fset.Position(start).Offset
	endPos := fset.Position(end).Offset

	// Validate positions
	if startPos < 0 || endPos < 0 || startPos > len(fileText) || endPos > len(fileText) {
		return "", fmt.Errorf("invalid positions")
	}

	// Extract the text within the specified range
	return fileText[startPos:endPos], nil
}
