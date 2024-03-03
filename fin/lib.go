package main

import "strings"

func nl2br(in string) string {
	return strings.Replace(in, "\n", "<br>\n", -1)
}
