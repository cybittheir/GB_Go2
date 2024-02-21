package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func writeFile(fn string) {

	f, err := os.Create(fn)
	if err != nil {
		fmt.Println("File not", err)
		return
	}
	fmt.Println("Введите предложение (exit или wr для выхода из режима заполнения):")
	reader := bufio.NewReader(os.Stdin)
	var n int
	var tx string
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		tm := time.Now()
		t := strings.TrimSpace(text)

		if strings.ToLower(t) == "exit" || strings.ToLower(t) == "wr" {
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

	fmt.Println("Информация сохранена")

	// setting ReadOnly mode after saving
	if err := os.Chmod(fn, 0444); err != nil {
		fmt.Println(err)
	}
	fmt.Println("ReadOnly mode - done")
	fmt.Println("====================")
}

func readFile(fn string) {
	s, err := os.Stat(fn)
	if err != nil {
		fmt.Println("File error: ", err)
		return
	}

	sz := s.Size()

	if sz == 0 {
		fmt.Println("File is empty")
		return
	}

	buf := make([]byte, sz)

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("File not opened: ", err)
		return
	}
	defer f.Close()

	//	if _, err := io.ReadFull(f, buf); err != nil {
	if _, err := f.Read(buf); err != nil {
		fmt.Println("Cannot read file: ", err)
		return
	}

	fmt.Printf("File %s, size=%d bytes:\n", fn, sz)
	fmt.Printf("%s\n", buf)

}

func main() {
	fn := "list.txt"

	writeFile(fn)
	readFile(fn)

	fmt.Println("Bye.")
}
