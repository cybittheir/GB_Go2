package main

import (
	"fmt"
	"os"
)

func main() {
	fn := "list.txt"

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
	fmt.Println("Bye.")
}
