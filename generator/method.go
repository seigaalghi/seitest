package generator

import (
	"fmt"

	"github.com/seigaalghi/seitest/utils"
)

func MethodTestGenerator(f utils.Function, lines []string, structs []utils.Struct) ([]string, error) {
	fmt.Println(f.Recv)
	return lines, nil
}
