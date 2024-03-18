package text

import (
	"regexp"
	"strings"
)

func CamelToSnake(camel string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(camel, "${1}_${2}")

	snake = strings.ToLower(snake)

	return snake
}