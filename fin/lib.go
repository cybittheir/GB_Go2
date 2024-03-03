package main

import (
	"log"
	"strings"
)

func nl2br(in string) string {
	return strings.Replace(in, "\n", "<br>\n", -1)
}

func checkPanic(e error) {
	if e != nil {
		panic(e)
	}
}

func check(e error) {
	if e != nil {
		log.Println("Error:", e)
	}
}
