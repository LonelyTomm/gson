package gson

import (
	"fmt"
)

func Encode(source interface{}) string {
	switch source := source.(type) {
	case map[string]interface{}:
		return encodeMap(source)
	case []interface{}:
		return encodeSlice(source)
	}

	return ""
}

func encodeMap(source map[string]interface{}) string {
	result := "{"

	currentIndex := 0
	for key, value := range source {
		result += encodeString(key) + ":"
		result += encodeValue(value)

		if currentIndex < len(source)-1 {
			result += ","
		}

		currentIndex++
	}

	result += "}"
	return result
}

func encodeSlice(source []interface{}) string {
	result := "["
	for index, value := range source {
		result += encodeValue(value)

		if index < len(source)-1 {
			result += ","
		}
	}

	result += "]"
	return result
}

func encodeValue(value interface{}) string {
	var result string
	switch value := value.(type) {
	case float64:
		result = encodeFloat(value)
	case string:
		result = encodeString(value)
	case []interface{}:
		result = encodeSlice(value)
	case map[string]interface{}:
		result = encodeMap(value)
	case bool:
		result = encodeBoolean(value)
	default:
		result = "null"
	}

	return result
}

func encodeBoolean(value bool) string {
	if value {
		return "true"
	}

	return "false"
}

func encodeString(source string) string {
	return "\"" + source + "\""
}

func encodeFloat(value float64) string {
	return fmt.Sprintf("%f", value)
}
