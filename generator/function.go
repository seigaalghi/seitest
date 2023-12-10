package generator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/seigaalghi/seitest/utils"
)

func FuncTestGenerator(f utils.Function, lines []string) ([]string, error) {
	lines = append(lines, fmt.Sprintf(`func Test%s(t *testing.T){`, f.Name))
	lines = append(lines, parsePayloadToStruct(f.Payload))
	lines = append(lines, parseResponseToAssert(f.Result))
	lines = append(lines, parseTestRunner(f.Payload, f.Result, f.Name))
	lines = append(lines, "\n}")
	return lines, nil
}

func parsePayloadToStruct(payload string) string {
	payloadStruct := fmt.Sprintf("type payload struct {\n%s}", payloadToNamedVariable(payload))

	return payloadStruct
}

func payloadToNamedVariable(payload string) string {
	pieces := strings.Split(removeSuffixPrefixParentheses(payload), ", ")

	occupied := map[string]string{}
	var outputText string
	for _, piece := range pieces {
		fields := strings.Fields(piece)
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
		dataType := r
		if splitted := strings.Split(r, " "); len(splitted) == 2 {
			dataType = splitted[1]
		}
		responseArr[i] = fmt.Sprintf("%s %s", indexToArg(i), dataType)
	}
	return fmt.Sprintf(`tests := []struct {
		scenario string
		payload  payload
		assert func(t *testing.T, %s)
	}{
		// Put Your Scenario Here
	}`, strings.Join(responseArr, ", "))
}

func indexToArg(i int) string {
	return fmt.Sprintf("res%d", i)
}

func parseTestRunner(payload, response, functionName string) string {
	responseArr := strings.Split(removeSuffixPrefixParentheses(response), ",")
	for i := range responseArr {
		responseArr[i] = indexToArg(i)
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
	countMap := make(map[string]int)

	result := make([]string, len(arr))

	for i, element := range arr {
		if countMap[element] == 0 {
			result[i] = element
		} else {
			result[i] = fmt.Sprintf("%s%d", element, countMap[element])
		}
		countMap[element]++
	}

	return result
}
