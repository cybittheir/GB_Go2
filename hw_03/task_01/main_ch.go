package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func squareInt(a int) chan int {
	sCh := make(chan int)
	go func() {
		sCh <- a * a
		fmt.Print("a ^ 2 =")
	}()
	return sCh
}

func dInt(sCh chan int) chan int {
	tCh := make(chan int)
	go func() {
		a, ok := <-sCh
		if ok {
			fmt.Println(a)
			r := 2 * a
			tCh <- r
			fmt.Print("2 x X = ")
		}
	}()
	return tCh
}

func toNum(a string) int {
	r, err := strconv.Atoi(a)
	if err != nil {
		fmt.Println("Ошибка ввода. Требуется число")
		r = 0
	}
	fmt.Println("input = ", r)
	return r
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Введите целое число или 'стоп' для выхода: ")
		in, _ := reader.ReadString('\n')
		t := strings.TrimSpace(in)
		a := 0
		if strings.ToLower(t) == "стоп" || strings.ToLower(t) == "stop" {
			break
		} else {
			a = toNum(t)

		}
		fc := squareInt(a)
		sc := dInt(fc)
		fmt.Println(<-sc)
	}
}
