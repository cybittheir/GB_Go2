package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// формирование элемента
func buildItem(n int) (string, bool) {
	l := [2]string{"(", ")"}
	r := l[0]
	z := 1
	op := 1
	for i := 1; i < 2*n-1; i++ {
		m := 1
		if i <= n {
			if z > 0 && z < n {
				m = rand.Intn(2)
			} else if z == 0 {
				m = 0
			} else if z < 0 {
				fmt.Print("*", z)
				break
			}
		} else {
			if op == n {
				m = 1
			} else if z == 0 {
				m = 0
			} else if z > 0 {
				m = rand.Intn(2)
			} else if z < 0 || z > n {
				break
			}
		}
		if m == 0 {
			z++
			op++
		} else {
			z--
		}

		r += l[m]
	}
	r += l[1]
	z--
	if z != 0 || len(r) < 2*n { // на случай если не получилось собрать
		fmt.Println("failed")
		return "", false
	}
	return r, true

}

// сбор элементов в мапу для последующей передачи в результирующий массив
func combine(n int) []string {

	s := make(map[string]string, n*2-1)
	sr := make([]string, 0, n)
	var f bool
	r := "()"
	w := catalana(n)
	fmt.Println("waiting for... ", w, "results")
	rand.Seed(time.Now().UnixNano())

	if n == 1 {
		s[r] = r
	} else {
		for {

			r, f = buildItem(n)
			if f == true && s[r] != r {
				s[r] = r
				if len(s) == w {
					break
				}
			}
		}
	}
	for i, _ := range s {
		sr = append(sr, i)
	}
	fmt.Println("===========")
	fmt.Printf("Вариантов: %d\n", len(s))

	return sr
}

// вычисление количества комбинаций
func catalana(i int) int {
	if i == 2 {
		return 2
	} else if i == 1 {
		return 1
	} else {
		return catalana(i-1) * (4*i - 2) / (i + 1)
	}
}

func main() {

	fmt.Print("Введите количество пар: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	n, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		fmt.Printf("Ошибка ввода. Требуется цифра. По-умолчанию - 1, ввели %v\n", n)
		n = 1
	}

	r := combine(n)
	fmt.Print("[")
	for _, v := range r {
		fmt.Printf("\"%s\"", v)
	}
	fmt.Print("]\n")

}
