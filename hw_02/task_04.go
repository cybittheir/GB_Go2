package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func writeFile(fn string) {

	s, err := os.Stat(fn)
	if errors.Is(err, os.ErrNotExist) {

	} else if err != nil {
		fmt.Println("File error: ", err)
		return
	} else if s.Mode() == 0444 || s.Mode() == 0440 || s.Mode() == 0400 {
		fmt.Println("Внимание! Файл не может быть открыт для записи!")
		return
	}

	var b bytes.Buffer

	fmt.Println("Введите предложение (exit или wr для выхода из режима заполнения):")
	reader := bufio.NewReader(os.Stdin)
	var n int

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
		b.WriteString(tt)
	}

	if err := os.WriteFile(fn, b.Bytes(), 0600); err != nil {
		fmt.Println("Unable to save", err)
		return
	}
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

	//	buf := make([]byte, sz)

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("Can not open file: ", err)
		return
	}
	defer f.Close()

	if rbytes, err := io.ReadAll(f); err != nil {
		fmt.Println("Cannot read file: ", err)
		return
	} else {
		fmt.Printf("File %s, size=%d bytes:\n", fn, sz)
		fmt.Printf("%s\n", string(rbytes))
	}

}

func main() {
	fn := "list.txt"

	writeFile(fn)
	readFile(fn)

	fmt.Println("Bye.")
}
