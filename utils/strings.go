package utlis

import "strings"

// remove ```json from start and \n``` from end of the string
// remove escape characters from the string

func FixJsonString(jsonString string) string {

	jsonString = strings.TrimPrefix(jsonString, "```json\n")
	jsonString = strings.TrimSuffix(jsonString, "```")

	jsonString = strings.ReplaceAll(jsonString, "\\\"", "\"")
	jsonString = strings.ReplaceAll(jsonString, "\\n", "\n")
	jsonString = strings.ReplaceAll(jsonString, "\\t", "\t")

	return jsonString
}
