package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
Задание 1. Конвейер

Цели задания:
-Научиться работать с каналами и горутинами.
-Понять, как должно происходить общение между потоками.

Что нужно сделать:
-Реализуйте паттерн-конвейер:
Программа принимает числа из стандартного ввода в бесконечном цикле и передаёт число в горутину.
-Квадрат: горутина высчитывает квадрат этого числа и передаёт в следующую горутину.
-Произведение: следующая горутина умножает квадрат числа на 2.
-При вводе «стоп» выполнение программы останавливается.

Советы и рекомендации:
Воспользуйтесь небуферизированными каналами и waitgroup.
*/

func sqInt(a int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	sCh := make(chan int)

	go func() {
		sCh <- a * a
		fmt.Print("X ^ 2 = Y = ")
	}()
	return sCh
}

func dbInt(sCh chan int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	tCh := make(chan int)
	go func() {
		a, ok := <-sCh
		if ok {
			fmt.Println(a)
			r := 2 * a
			tCh <- r
			fmt.Print("2 x Y = ")
		}
	}()
	return tCh
}

func intoNum(a string) int {
	r, err := strconv.Atoi(a)
	if err != nil {
		fmt.Println("Ошибка ввода. Требуется число")
		r = 0
	}
	fmt.Println("input = X = ", r)
	return r
}

func greeting() (int, bool) {
	reader := bufio.NewReader(os.Stdin)

	ok := false
	num := 0
	fmt.Print("Введите целое число или 'стоп' для выхода: ")
	in, _ := reader.ReadString('\n')
	t := strings.TrimSpace(in)
	if strings.ToLower(t) == "стоп" || strings.ToLower(t) == "stop" {
		ok = false
	} else {
		num = intoNum(t)
		ok = true
	}

	return num, ok
}

func main() {

	var wg sync.WaitGroup
	for {
		in, ok := greeting()
		if !ok {
			break
		} else if in > 0 {
			wg.Add(2)
			fc := sqInt(in, &wg)
			sc := dbInt(fc, &wg)
			fmt.Println(<-sc)
			wg.Wait()
		}
	}
	fmt.Println("Bye")
}
