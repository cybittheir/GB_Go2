package main

/*
Задание 2. Сортировка пузырьком
Что нужно сделать:
Отсортируйте массив длиной шесть пузырьком.
*/

import "fmt"

func BubbleSort(a *[6]int) {

	var t int
	var w bool
	//	s := *a
	for {
		w = true
		for i := range 5 {
			if a[i+1] < a[i] {
				t = a[i]
				a[i] = a[i+1]
				a[i+1] = t
				w = false
			}
		}
		if w {
			break
		}
	}
	fmt.Println(a)

}

func main() {
	a := [6]int{43, 10, 35, 14, 17, 22}

	fmt.Println("Source: ", a)

	BubbleSort(&a)

	fmt.Println("Result: ", a)

}
