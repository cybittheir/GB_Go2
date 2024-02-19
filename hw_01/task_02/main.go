package main

/*
Задание 2. Сортировка пузырьком
Что нужно сделать:
Отсортируйте массив длиной шесть пузырьком.
*/

import "fmt"

func BubbleSort(a *[]int) {

	var t int
	var w bool
	s := *a
	for {
		w = true
		for i := range len(s) - 1 {
			if s[i+1] < s[i] {
				t = s[i]
				s[i] = s[i+1]
				s[i+1] = t
				w = false
			}
		}
		if w {
			break
		}
	}
	fmt.Println(s)

}

func main() {
	a := []int{43, 10, 35, 14, 17, 22}

	fmt.Println(a)

	BubbleSort(&a)

	fmt.Println(a)

}
