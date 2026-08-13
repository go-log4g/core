package impl

import "fmt"

type ParameterizedMessageFormatter struct {
}

func NewParameterizedMessageFormatter() *ParameterizedMessageFormatter {
	return &ParameterizedMessageFormatter{}
}

func (this *ParameterizedMessageFormatter) Format(pattern string, args ...any) string {
	if len(args) == 0 {
		return pattern
	}

	result := make([]byte, 0, len(pattern)+32)
	argIndex := 0

	for index := 0; index < len(pattern); {
		if index+1 >= len(pattern) || pattern[index] != '{' || pattern[index+1] != '}' {
			result = append(result, pattern[index])
			index++
			continue
		}

		backslashes := 0
		for i := index - 1; i >= 0 && pattern[i] == '\\'; i-- {
			backslashes++
		}

		if backslashes > 0 {
			result = result[:len(result)-backslashes]
			for range backslashes / 2 {
				result = append(result, '\\')
			}
		}

		if backslashes%2 == 1 {
			result = append(result, '{', '}')
			index += 2
			continue
		}

		if argIndex >= len(args) {
			result = append(result, '{', '}')
			index += 2
			continue
		}

		result = fmt.Append(result, args[argIndex])
		argIndex++
		index += 2
	}

	return string(result)
}
