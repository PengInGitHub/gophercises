package main

import (
	"bytes"
	"regexp"
)

func main() {
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch <= '9' && ch >= '0' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func normalizeREGEXP(phone string) string {
	re := regexp.MustCompile("\\D") //find all non int
	return re.ReplaceAllString(phone, "")
}
