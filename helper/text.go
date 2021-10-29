package helper

import "regexp"

var regex *regexp.Regexp

func init() {
	regex, _ = regexp.Compile("\\t|\\n")
}

func EscapeNewLineAndIndent(text string) string {
	return regex.ReplaceAllString(text, "")
}
