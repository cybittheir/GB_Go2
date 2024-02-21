package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fn := "list.txt"
	f, err := os.Create(fn)
	if err != nil {
		fmt.Println("File not", err)
		return
	}
	fmt.Println("Введите предложение:")
	reader := bufio.NewReader(os.Stdin)
	var n int
	var tx string
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		tm := time.Now()
		t := strings.TrimSpace(text)

		if strings.ToLower(t) == "exit" {
			break
		} else if t == "" {
			continue
		}
		n++
		tt := fmt.Sprintf("%d. %s %s\n", n, tm.Format("02/01/2006 15:04:05"), t)
		fmt.Println(tt)
		tx += tt
	}
	defer f.Close()

	f.WriteString(tx)

	// setting ReadOnly mode after saving
	if err := os.Chmod(fn, 0444); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Bye.")
}
