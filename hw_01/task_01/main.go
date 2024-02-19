package main

import (
	"fmt"
)

/*
Задание 1. Слияние отсортированных массивов
Что нужно сделать:
Напишите функцию, которая производит слияние двух отсортированных массивов длиной четыре и пять в один массив длиной девять.
*/

func splitSorted(a1 []int, a2 []int) []int {

	var r []int

	j, k := 0, 0   // slices indexes
	sj, sk := 0, 0 // flags switched when last element in slice was used

	for i := range len(a1) + len(a2) {
		switch {
		case j == len(a1) || k == len(a2):
			if k == len(a2) {
				k--
				r = append(r, a1[j])
				fmt.Println("ob: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a1[j])
			} else {
				j--
				r = append(r, a2[k])
				fmt.Println("0a: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a2[k])
			}
		case j == len(a1)-1:
			if a2[k] <= a1[j] || sj == 1 {
				r = append(r, a2[k])
				fmt.Println("1a: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a2[k])
				if k < len(a2) {
					k++
				}
			} else {
				r = append(r, a1[j])
				fmt.Println("1b: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a1[j])
				sj = 1
			}
		case k == len(a2)-1:
			if a1[j] <= a2[k] || sk == 1 {
				r = append(r, a1[j])
				fmt.Println("2a: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a1[j])
				if j < len(a1) {
					j++
				}
			} else {
				r = append(r, a2[k])
				fmt.Println("2b: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a2[k])
				sk = 1
			}
		case a1[j] > a2[k] && j < len(a1):
			r = append(r, a2[k])
			fmt.Println("3: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a2[k])
			if k < len(a2) {
				k++
			}
		case a2[k] >= a1[j] && k < len(a2):
			r = append(r, a1[j])
			fmt.Println("4: ", k, "=", a2[k], ";", j, "=", a1[j], ">", a1[j])
			if j < len(a1) {
				j++
			}
		}
		fmt.Println(i)
	}

	return r
}

func main() {
	a1 := [4]int{1, 3, 9, 25}
	a2 := [5]int{10, 11, 14, 22, 23}

	fmt.Println(a1, " + ", a2)

	s := splitSorted(a1[0:], a2[0:])

	var r [9]int

	copy(r[:], s[:]) // transfer from slice to array

	fmt.Println(r)

}
