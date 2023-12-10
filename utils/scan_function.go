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
	Recv     string
}

type Field struct {
	Name     string
	DataType string
}

type Struct struct {
	FilePath string
	Name     string
	Content  string
	Fields   []Field
}

type Import struct {
	FilePath string
	Content  string
}

func ScanFunctions(packagePath string) ([]Function, []Struct, []Import, error) {
	var functions []Function
	var structs []Struct
	var imports []Import
	fset := token.NewFileSet()

	err := filepath.Walk(packagePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go") {
			node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return err
			}
			for _, decl := range node.Decls {
				if n, ok := decl.(*ast.GenDecl); ok {
					if n.Tok.String() == "type" {
						for _, spec := range n.Specs {
							if ts, ok := spec.(*ast.TypeSpec); ok {
								if st, ok := ts.Type.(*ast.StructType); ok {
									structName := ts.Name.Name
									content, err := getFileText(path, st.Pos(), st.End(), fset)
									if err != nil {
										log.Fatal(err.Error())
									}

									var fields []Field
									for _, field := range st.Fields.List {
										fieldName := field.Names[0].Name
										fieldType, err := getFileText(path, field.Type.Pos(), field.Type.End(), fset)
										if err != nil {
											log.Fatal(err.Error())
										}
										fields = append(fields, Field{
											Name:     fieldName,
											DataType: fieldType,
										})
									}

									structs = append(structs, Struct{
										FilePath: path,
										Content:  content,
										Name:     structName,
										Fields:   fields,
									})
								}
							}
						}
					}

					if n.Tok.String() == "import" {
						importText, err := getFileText(path, n.Pos(), n.End(), fset)
						if err != nil {
							log.Fatal(err.Error())
						}

						imports = append(imports, Import{
							FilePath: path,
							Content:  importText,
						})
					}
				}
				if fn, ok := decl.(*ast.FuncDecl); ok {
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

					var recv string
					if fn.Recv != nil {
						recv, err = getFileText(path, fn.Recv.Pos(), fn.Recv.End(), fset)
						if err != nil {
							log.Fatal(err.Error())
						}
					}

					functions = append(functions, Function{
						FilePath: path,
						Name:     fn.Name.Name,
						Content:  funcText,
						IsMethod: fn.Recv != nil,
						Package:  node.Name.Name,
						Payload:  payload,
						Result:   result,
						Recv:     recv,
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return functions, structs, imports, nil
}

func getFileText(filePath string, start, end token.Pos, fset *token.FileSet) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fileText := string(data)
	startPos := fset.Position(start).Offset
	endPos := fset.Position(end).Offset

	if startPos < 0 || endPos < 0 || startPos > len(fileText) || endPos > len(fileText) {
		return "", fmt.Errorf("invalid positions")
	}

	return fileText[startPos:endPos], nil
}
