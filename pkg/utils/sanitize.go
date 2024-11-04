package utils

import (
	"github.com/microcosm-cc/bluemonday"
)

func Sanitize(str string) string {
	sanitizer := bluemonday.UGCPolicy()
	return sanitizer.Sanitize(str)
}