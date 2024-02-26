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

func square(wg *sync.WaitGroup, a int) int {
	defer wg.Done()
	r := a * a
	fmt.Println("X^2 = ", r)
	go double(wg, r)
	return r
}

func double(wg *sync.WaitGroup, a int) int {
	defer wg.Done()
	r := 2 * a
	fmt.Println("2 x X = ", r)
	return r
}

func inNum(a string) int {
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
	var wg sync.WaitGroup
	//	type(ch1 chan)
	for {
		fmt.Print("Введите целое число или 'стоп' для выхода: ")
		in, _ := reader.ReadString('\n')
		t := strings.TrimSpace(in)
		a := 0
		if strings.ToLower(t) == "стоп" || strings.ToLower(t) == "stop" {
			break
		} else {
			a = inNum(t)

		}

		wg.Add(2)
		go square(&wg, a)

		wg.Wait()
	}
}
